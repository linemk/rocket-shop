package main

import (
	"context"
	"log"
	"net"

	"github.com/google/uuid"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = "50052"
)

// PaymentService реализует gRPC сервис для работы с платежами
type PaymentService struct {
	payment_v1.UnimplementedPaymentServiceServer
}

// NewPaymentService создает новый экземпляр PaymentService
func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// PayOrder обрабатывает команду на оплату и возвращает transaction_uuid
func (s *PaymentService) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	// Генерируем UUID транзакции
	transactionUUID := uuid.New()

	// Выводим сообщение в консоль согласно спецификации
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID.String())

	// Возвращаем ответ
	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}

func main() {
	// Создаем gRPC сервер
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Создаем и регистрируем PaymentService
	paymentService := NewPaymentService()
	payment_v1.RegisterPaymentServiceServer(grpcServer, paymentService)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	log.Printf("💳 PaymentService starting on port %s", grpcPort)

	// Запускаем сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
