package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/linemk/rocket-shop/payment/internal/delivery/v1"
	paymentRepository "github.com/linemk/rocket-shop/payment/internal/repository/payment"
	"github.com/linemk/rocket-shop/payment/internal/usecase"
	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
)

const (
	grpcPort = "50052"
)

func main() {
	// Создаем репозиторий
	paymentRepo := paymentRepository.NewRepository()

	// Создаем UseCase
	paymentUseCase := usecase.NewUseCase(paymentRepo)

	// Создаем API handler
	api := v1.NewAPI(paymentUseCase)

	// Создаем gRPC сервер
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Регистрируем PaymentService
	payment_v1.RegisterPaymentServiceServer(grpcServer, api)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	log.Printf("💳 PaymentService starting on port %s", grpcPort)

	// Запускаем сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
