//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

// InsertTestPart — вставляет тестовую деталь в коллекцию Mongo и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"_id":            partUUID,
		"uuid":           partUUID,
		"name":           gofakeit.CarModel(),
		"description":    gofakeit.Sentence(10),
		"price":          gofakeit.Float64Range(100, 10000),
		"stock_quantity": int64(gofakeit.IntRange(0, 1000)),
		"category":       1, // CATEGORY_ENGINE
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(10, 200),
			"width":  gofakeit.Float64Range(10, 200),
			"height": gofakeit.Float64Range(10, 200),
			"weight": gofakeit.Float64Range(1, 500),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags":       []string{gofakeit.Word(), gofakeit.Word()},
		"metadata":   bson.M{"test": "data"},
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-test"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// InsertTestPartWithData — вставляет тестовую деталь с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part *inventoryV1.Part) (string, error) {
	partUUID := part.GetUuid()
	if partUUID == "" {
		partUUID = gofakeit.UUID()
	}
	now := time.Now()

	metadata := bson.M{}
	if part.GetMetadata() != nil {
		metadata = part.GetMetadata().AsMap()
	}

	partDoc := bson.M{
		"_id":            partUUID,
		"uuid":           partUUID,
		"name":           part.GetName(),
		"description":    part.GetDescription(),
		"price":          part.GetPrice(),
		"stock_quantity": part.GetStockQuantity(),
		"category":       int32(part.GetCategory()),
		"dimensions": bson.M{
			"length": part.GetDimensions().GetLength(),
			"width":  part.GetDimensions().GetWidth(),
			"height": part.GetDimensions().GetHeight(),
			"weight": part.GetDimensions().GetWeight(),
		},
		"manufacturer": bson.M{
			"name":    part.GetManufacturer().GetName(),
			"country": part.GetManufacturer().GetCountry(),
			"website": part.GetManufacturer().GetWebsite(),
		},
		"tags":       part.GetTags(),
		"metadata":   metadata,
		"created_at": primitive.NewDateTimeFromTime(part.GetCreatedAt().AsTime()),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-test"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// GetTestPart — возвращает тестовую деталь
func (env *TestEnvironment) GetTestPart() *inventoryV1.Part {
	metadata, _ := structpb.NewStruct(map[string]interface{}{
		"test_field": "test_value",
	})

	return &inventoryV1.Part{
		Uuid:          gofakeit.UUID(),
		Name:          "Quantum Engine X-3000",
		Description:   "High-performance quantum propulsion engine",
		Price:         15999.99,
		StockQuantity: 50,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 150.0,
			Width:  100.0,
			Height: 80.0,
			Weight: 500.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "SpaceTech Industries",
			Country: "USA",
			Website: "https://spacetech.com",
		},
		Tags:      []string{"quantum", "engine", "propulsion"},
		Metadata:  metadata,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-test"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
