package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	inventoryClient "github.com/linemk/rocket-shop/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/linemk/rocket-shop/order/internal/client/grpc/payment/v1"
	v1 "github.com/linemk/rocket-shop/order/internal/delivery/v1"
	"github.com/linemk/rocket-shop/order/internal/repository"
	"github.com/linemk/rocket-shop/order/internal/usecase"
	order_v1 "github.com/linemk/rocket-shop/shared/pkg/openapi/order/v1"
)

const (
	httpPort             = "8080"
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 10 * time.Second
)

func main() {
	// Инициализируем репозиторий
	orderRepository := repository.NewRepository()

	// Инициализируем клиенты
	inventoryClient, err := inventoryClient.NewClient(inventoryServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create inventory client: %v", err)
	}

	paymentClient, err := paymentClient.NewClient(paymentServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create payment client: %v", err)
	}

	// Инициализируем UseCase
	orderUseCase := usecase.NewUseCase(orderRepository, inventoryClient, paymentClient)

	// Инициализируем API
	api := v1.NewAPI(orderUseCase)

	// Создаем сервер
	orderServer, err := order_v1.NewServer(api)
	if err != nil {
		log.Fatalf("Failed to create order server: %v", err)
	}

	// Устанавливаем defer после всех проверок ошибок
	defer func() {
		if err := inventoryClient.Close(); err != nil {
			log.Printf("Failed to close inventory client: %v", err)
		}
		if err := paymentClient.Close(); err != nil {
			log.Printf("Failed to close payment client: %v", err)
		}
	}()

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

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
