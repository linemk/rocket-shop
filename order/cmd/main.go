package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	// Адрес InventoryService
	inventoryServiceAddr = "localhost:50051"
	// Адрес PaymentService
	paymentServiceAddr = "localhost:50052"
)

type OrderService struct {
	mu              sync.RWMutex
	orders          map[string]*Order
	inventoryClient inventory_v1.InventoryServiceClient
	inventoryConn   *grpc.ClientConn
	paymentClient   payment_v1.PaymentServiceClient
	paymentConn     *grpc.ClientConn
}

type Order struct {
	ID            string                 `json:"order_id"`
	UserID        string                 `json:"user_id"`
	DetailsID     []string               `json:"details_id"`
	TotalPrice    float32                `json:"total_price"`
	TransactionID string                 `json:"transaction_id"`
	PaymentMethod order_v1.PaymentMethod `json:"payment_method"`
	Status        order_v1.OrderStatus   `json:"status"`
}

func NewOrderService() (*OrderService, error) {
	// Подключаемся к InventoryService
	inventoryConn, err := grpc.Dial(inventoryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to InventoryService: %w", err)
	}

	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryConn)

	// Подключаемся к PaymentService
	paymentConn, err := grpc.Dial(paymentServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		inventoryConn.Close()
		return nil, fmt.Errorf("failed to connect to PaymentService: %w", err)
	}

	paymentClient := payment_v1.NewPaymentServiceClient(paymentConn)

	return &OrderService{
		orders:          make(map[string]*Order),
		inventoryClient: inventoryClient,
		inventoryConn:   inventoryConn,
		paymentClient:   paymentClient,
		paymentConn:     paymentConn,
	}, nil
}

func (s *OrderService) GetOrderByUUID(id string) (*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("order by id %s not found", id)
	}

	return w, nil
}

// validateParts проверяет существование деталей через InventoryService и возвращает их общую стоимость
func (s *OrderService) validateParts(ctx context.Context, partUUIDs []uuid.UUID) (float32, error) {
	if len(partUUIDs) == 0 {
		return 0, fmt.Errorf("no parts specified")
	}

	var totalPrice float32

	for _, partUUID := range partUUIDs {
		// Получаем информацию о детали через InventoryService
		resp, err := s.inventoryClient.GetPart(ctx, &inventory_v1.GetPartRequest{
			Uuid: partUUID.String(),
		})
		if err != nil {
			return 0, fmt.Errorf("part %s not found in inventory: %w", partUUID.String(), err)
		}

		// Проверяем наличие на складе
		if resp.Part.StockQuantity <= 0 {
			return 0, fmt.Errorf("part %s is out of stock", partUUID.String())
		}

		// Добавляем цену к общей стоимости
		totalPrice += float32(resp.Part.Price)
	}

	return totalPrice, nil
}

func (s *OrderService) GetOrder(_ context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	orderUUID := params.OrderUUID.String()

	order, err := s.GetOrderByUUID(orderUUID)
	if err != nil || order == nil {
		return &order_v1.NotFoundErr{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	var transactionUUID uuid.UUID
	if order.TransactionID != "" {
		transactionUUID = uuid.MustParse(order.TransactionID)
	}

	response := &order_v1.GetOrderResp{
		OrderUUID:       params.OrderUUID, // или uuid.MustParse(order.ID)
		UserUUID:        uuid.MustParse(order.UserID),
		PartUuids:       convertStringSliceToUUIDSlice(order.DetailsID),
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   order_v1.PaymentMethod(order.PaymentMethod),
		Status:          order_v1.OrderStatus(order.Status),
	}

	return response, nil
}

func (s *OrderService) NewError(_ context.Context, err error) *order_v1.UnexpectedErrStatusCode {
	return &order_v1.UnexpectedErrStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.UnexpectedErr{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func (s *OrderService) PayOrder(ctx context.Context, req *order_v1.PayOrderReq, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	orderID := params.OrderUUID.String()
	paymentMethod := order_v1.PaymentMethod(req.PaymentMethod)
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return &order_v1.NotFoundErr{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	if order.Status != order_v1.OrderStatusPENDINGPAYMENT {
		return &order_v1.ConflictErr{
			Code:    409,
			Message: "Order cannot be paid in current status",
		}, nil
	}

	// Вызываем PaymentService для обработки платежа
	paymentResp, err := s.paymentClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderID,
		UserUuid:      order.UserID,
		PaymentMethod: convertOpenAPIPaymentMethodToProto(paymentMethod),
	})
	if err != nil {
		return &order_v1.BadRequest{
			Code:    400,
			Message: fmt.Sprintf("Payment failed: %v", err),
		}, nil
	}

	// Обновляем заказ
	order.PaymentMethod = paymentMethod
	order.TransactionID = paymentResp.TransactionUuid
	order.Status = order_v1.OrderStatusPAID

	// Конвертируем UUID для ответа
	transactionUUID, err := uuid.Parse(paymentResp.TransactionUuid)
	if err != nil {
		return &order_v1.BadRequest{
			Code:    400,
			Message: "Invalid transaction UUID",
		}, nil
	}

	return &order_v1.PayOrderResp{
		TransactionUUID: transactionUUID,
	}, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req order_v1.OptCreateOrderReq) (order_v1.CreateOrderRes, error) {
	userID := req.Value.UserUUID.String()
	partsIDs := make([]string, 0, len(req.Value.PartUuids))
	for _, id := range req.Value.PartUuids {
		partsIDs = append(partsIDs, id.String())
	}

	// Проверяем детали через InventoryService и получаем общую стоимость
	totalPrice, err := s.validateParts(ctx, req.Value.PartUuids)
	if err != nil {
		return &order_v1.BadRequest{
			Code:    400,
			Message: fmt.Sprintf("Invalid parts: %v", err),
		}, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	orderUUID := uuid.New()
	order := &Order{
		ID:         orderUUID.String(),
		UserID:     userID,
		DetailsID:  partsIDs,
		TotalPrice: totalPrice,
		Status:     order_v1.OrderStatusPENDINGPAYMENT,
	}

	s.orders[order.ID] = order

	return &order_v1.CreateOrderResp{
		UUID:       orderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	orderID := params.OrderUUID.String()
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return &order_v1.NotFoundErr{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	if order.Status == order_v1.OrderStatusPENDINGPAYMENT {
		order.Status = order_v1.OrderStatusCANCELLED
		return &order_v1.CancelOrderNoContent{}, nil
	}

	if order.Status == order_v1.OrderStatusPAID {
		return &order_v1.ConflictErr{
			Code:    409,
			Message: "Order already paid and cannot be cancelled",
		}, nil
	}

	return &order_v1.CancelOrderNoContent{}, nil
}

func convertStringSliceToUUIDSlice(strSlice []string) []uuid.UUID {
	uuidSlice := make([]uuid.UUID, len(strSlice))
	for i, str := range strSlice {
		uuidSlice[i] = uuid.MustParse(str)
	}
	return uuidSlice
}

// convertOpenAPIPaymentMethodToProto конвертирует PaymentMethod из OpenAPI в proto
func convertOpenAPIPaymentMethodToProto(openAPIMethod order_v1.PaymentMethod) payment_v1.PaymentMethod {
	switch openAPIMethod {
	case order_v1.PaymentMethodPAYMENTMETHODUNSPECIFIED:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	case order_v1.PaymentMethodPAYMENTMETHODCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case order_v1.PaymentMethodPAYMENTMETHODSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case order_v1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case order_v1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func main() {
	s, err := NewOrderService()
	if err != nil {
		log.Fatalf("Failed to create order service: %v", err)
	}
	defer func() {
		err = s.inventoryConn.Close()
		if err != nil {
			log.Fatalf("Failed to close inventory connection: %v", err)
		}
		err = s.paymentConn.Close()
		if err != nil {
			log.Fatalf("Failed to close payment connection: %v", err)
		}
	}()

	orderServer, err := order_v1.NewServer(s)
	if err != nil {
		log.Fatalf("Failed to create order server: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		log.Printf("🔗 Подключен к InventoryService на %s\n", inventoryServiceAddr)
		log.Printf("💳 Подключен к PaymentService на %s\n", paymentServiceAddr)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
