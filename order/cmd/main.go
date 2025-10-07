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

	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderService struct {
	mu     sync.RWMutex
	orders map[string]*Order
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

func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]*Order),
	}
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

func (s *OrderService) PayOrder(_ context.Context, req *order_v1.PayOrderReq, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
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

	// TODO: обращение к PaymentService через gRPC
	transactionUUID := uuid.New()

	// Обновляем заказ
	order.PaymentMethod = paymentMethod
	order.TransactionID = transactionUUID.String()
	order.Status = order_v1.OrderStatusPAID

	return &order_v1.PayOrderResp{
		TransactionUUID: transactionUUID,
	}, nil
}

func (s *OrderService) CreateOrder(_ context.Context, req order_v1.OptCreateOrderReq) (order_v1.CreateOrderRes, error) {
	userID := req.Value.UserUUID.String()
	partsIDs := make([]string, 0, len(req.Value.PartUuids))
	for _, id := range req.Value.PartUuids {
		partsIDs = append(partsIDs, id.String())
	}

	// TODO: проверка деталей через InventoryService
	// TODO: расчет цены через InventoryService

	s.mu.Lock()
	defer s.mu.Unlock()

	orderUUID := uuid.New()
	order := &Order{
		ID:         orderUUID.String(),
		UserID:     userID,
		DetailsID:  partsIDs,
		TotalPrice: 100.43, // TODO доделать при создании сервиса расчетов
		Status:     order_v1.OrderStatusPENDINGPAYMENT,
	}

	s.orders[order.ID] = order

	return &order_v1.CreateOrderResp{
		OrderUUID:  orderUUID,
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

func main() {
	s := NewOrderService()

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
