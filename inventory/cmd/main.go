package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/linemk/rocket-shop/inventory/internal/converter"
	v1 "github.com/linemk/rocket-shop/inventory/internal/delivery/v1"
	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	inventoryRepository "github.com/linemk/rocket-shop/inventory/internal/repository/inventory"
	"github.com/linemk/rocket-shop/inventory/internal/usecase"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort           = "50051"
	defaultMongoURI    = "mongodb://inventory_user:inventory_password@localhost:27017"
	defaultMongoDBName = "inventory_db"
)

func main() {
	ctx := context.Background()

	// Получаем параметры подключения к MongoDB из переменных окружения
	mongoURI := os.Getenv("INVENTORY_MONGO_URI")
	if mongoURI == "" {
		mongoURI = defaultMongoURI
	}

	mongoDBName := os.Getenv("INVENTORY_MONGO_DB")
	if mongoDBName == "" {
		mongoDBName = defaultMongoDBName
	}

	// Создаем клиент MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Проверяем соединение
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return
	}

	log.Println("Successfully connected to MongoDB")

	// Получаем базу данных
	db := client.Database(mongoDBName)

	// Создаем MongoDB репозиторий
	inventoryRepo := inventoryRepository.NewMongoRepository(db)

	// Инициализируем тестовые данные
	initTestData(inventoryRepo)

	// Создаем UseCase
	inventoryUseCase := usecase.NewUseCase(inventoryRepo)

	// Создаем API handler
	api := v1.NewAPI(inventoryUseCase)

	// Создаем gRPC сервер
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Printf("Failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()

	// Регистрируем InventoryService
	inventory_v1.RegisterInventoryServiceServer(grpcServer, api)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	log.Printf("🚀 InventoryService starting on port %s", grpcPort)

	// Запускаем сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("Failed to serve: %v", err)
	}
}

// initTestData инициализирует тестовые данные деталей
func initTestData(repo interface {
	GetPart(ctx context.Context, uuid string) (models.Part, error)
	CreatePart(ctx context.Context, part models.Part) error
},
) {
	now := time.Now()

	// Создаем тестовые детали
	testParts := []*inventory_v1.Part{
		{
			Uuid:          "123e4567-e89b-12d3-a456-426614174001",
			Name:          "Main Engine",
			Description:   "Powerful rocket engine for main propulsion",
			Price:         50000.0,
			StockQuantity: 5,
			Category:      inventory_v1.Category_CATEGORY_ENGINE,
			Dimensions: &inventory_v1.Dimensions{
				Length: 200.0,
				Width:  50.0,
				Height: 50.0,
				Weight: 1000.0,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "SpaceTech Corp",
				Country: "USA",
				Website: "https://spacetech.com",
			},
			Tags:      []string{"engine", "propulsion", "main"},
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
		{
			Uuid:          "123e4567-e89b-12d3-a456-426614174002",
			Name:          "Liquid Hydrogen Fuel",
			Description:   "High-efficiency fuel for rocket engines",
			Price:         1000.0,
			StockQuantity: 100,
			Category:      inventory_v1.Category_CATEGORY_FUEL,
			Dimensions: &inventory_v1.Dimensions{
				Length: 100.0,
				Width:  100.0,
				Height: 200.0,
				Weight: 500.0,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "Fuel Systems Inc",
				Country: "Germany",
				Website: "https://fuelsystems.de",
			},
			Tags:      []string{"fuel", "hydrogen", "liquid"},
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
		{
			Uuid:          "123e4567-e89b-12d3-a456-426614174003",
			Name:          "Observation Porthole",
			Description:   "Reinforced glass porthole for crew observation",
			Price:         5000.0,
			StockQuantity: 20,
			Category:      inventory_v1.Category_CATEGORY_PORTHOLE,
			Dimensions: &inventory_v1.Dimensions{
				Length: 30.0,
				Width:  30.0,
				Height: 5.0,
				Weight: 10.0,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "GlassTech Ltd",
				Country: "Japan",
				Website: "https://glasstech.jp",
			},
			Tags:      []string{"porthole", "glass", "observation"},
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
		{
			Uuid:          "123e4567-e89b-12d3-a456-426614174004",
			Name:          "Stabilizer Wing",
			Description:   "Aerodynamic wing for flight stabilization",
			Price:         15000.0,
			StockQuantity: 8,
			Category:      inventory_v1.Category_CATEGORY_WING,
			Dimensions: &inventory_v1.Dimensions{
				Length: 300.0,
				Width:  100.0,
				Height: 20.0,
				Weight: 200.0,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "AeroDynamics Corp",
				Country: "France",
				Website: "https://aerodynamics.fr",
			},
			Tags:      []string{"wing", "stabilizer", "aerodynamic"},
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
	}

	// Добавляем метаданные для некоторых деталей
	metadata := map[string]interface{}{
		"certification":     "ISO-9001",
		"warranty_years":    5,
		"temperature_range": "-200 to 2000°C",
	}

	structMetadata, err := structpb.NewStruct(metadata)
	if err != nil {
		log.Printf("Failed to create struct metadata: %v", err)
		return
	}
	testParts[0].Metadata = structMetadata
	testParts[1].Metadata = structMetadata

	// Добавляем детали в репозиторий
	for _, protoPart := range testParts {
		// Конвертируем protobuf в модель
		part := converter.ProtoToPart(protoPart)

		// Проверяем, существует ли деталь
		_, err := repo.GetPart(context.TODO(), part.UUID)
		if err == nil {
			// Деталь уже существует, пропускаем
			log.Printf("Part %s already exists, skipping", part.UUID)
			continue
		}

		// Деталь не существует, создаем её
		if err := repo.CreatePart(context.TODO(), part); err != nil {
			log.Printf("Failed to create part %s: %v", part.UUID, err)
		} else {
			log.Printf("Created part %s", part.UUID)
		}
	}
}
