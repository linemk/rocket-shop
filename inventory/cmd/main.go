package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort = "50051"
)

// InventoryService реализует gRPC сервис для работы с инвентарем деталей
type InventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*inventory_v1.Part
}

// NewInventoryService создает новый экземпляр InventoryService
func NewInventoryService() *InventoryService {
	service := &InventoryService{
		parts: make(map[string]*inventory_v1.Part),
	}

	// Инициализируем тестовые данные
	service.initTestData()

	return service
}

// initTestData инициализирует тестовые данные деталей
func (s *InventoryService) initTestData() {
	now := timestamppb.New(time.Now())

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
			CreatedAt: now,
			UpdatedAt: now,
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
			CreatedAt: now,
			UpdatedAt: now,
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
			CreatedAt: now,
			UpdatedAt: now,
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
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Добавляем метаданные для некоторых деталей
	metadata := map[string]interface{}{
		"certification":     "ISO-9001",
		"warranty_years":    5,
		"temperature_range": "-200 to 2000°C",
	}

	structMetadata, _ := structpb.NewStruct(metadata)
	testParts[0].Metadata = structMetadata
	testParts[1].Metadata = structMetadata

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, part := range testParts {
		s.parts[part.Uuid] = part
	}
}

// GetPart возвращает информацию о детали по её UUID
func (s *InventoryService) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, exists := s.parts[req.Uuid]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Part with UUID %s not found", req.Uuid)
	}

	return &inventory_v1.GetPartResponse{
		Part: part,
	}, nil
}

// ListParts возвращает список деталей с возможностью фильтрации
func (s *InventoryService) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*inventory_v1.Part

	// Если фильтр пустой, возвращаем все детали
	if req.Filter == nil || isEmptyFilter(req.Filter) {
		for _, part := range s.parts {
			result = append(result, part)
		}
		return &inventory_v1.ListPartsResponse{Parts: result}, nil
	}

	// Применяем фильтрацию поэтапно
	candidates := make(map[string]*inventory_v1.Part)

	// Копируем все детали как кандидаты
	for uuid, part := range s.parts {
		candidates[uuid] = part
	}

	// Фильтр по UUID
	if len(req.Filter.Uuids) > 0 {
		filtered := make(map[string]*inventory_v1.Part)
		for _, uuid := range req.Filter.Uuids {
			if part, exists := candidates[uuid]; exists {
				filtered[uuid] = part
			}
		}
		candidates = filtered
	}

	// Фильтр по именам
	if len(req.Filter.Names) > 0 {
		filtered := make(map[string]*inventory_v1.Part)
		for uuid, part := range candidates {
			for _, name := range req.Filter.Names {
				if part.Name == name {
					filtered[uuid] = part
					break
				}
			}
		}
		candidates = filtered
	}

	// Фильтр по категориям
	if len(req.Filter.Categories) > 0 {
		filtered := make(map[string]*inventory_v1.Part)
		for uuid, part := range candidates {
			for _, category := range req.Filter.Categories {
				if part.Category == category {
					filtered[uuid] = part
					break
				}
			}
		}
		candidates = filtered
	}

	// Фильтр по странам производителей
	if len(req.Filter.ManufacturerCountries) > 0 {
		filtered := make(map[string]*inventory_v1.Part)
		for uuid, part := range candidates {
			for _, country := range req.Filter.ManufacturerCountries {
				if part.Manufacturer != nil && part.Manufacturer.Country == country {
					filtered[uuid] = part
					break
				}
			}
		}
		candidates = filtered
	}

	// Фильтр по тегам
	if len(req.Filter.Tags) > 0 {
		filtered := make(map[string]*inventory_v1.Part)
		for uuid, part := range candidates {
			for _, filterTag := range req.Filter.Tags {
				for _, partTag := range part.Tags {
					if partTag == filterTag {
						filtered[uuid] = part
						break
					}
				}
				if _, exists := filtered[uuid]; exists {
					break
				}
			}
		}
		candidates = filtered
	}

	// Преобразуем результат в слайс
	for _, part := range candidates {
		result = append(result, part)
	}

	return &inventory_v1.ListPartsResponse{Parts: result}, nil
}

// isEmptyFilter проверяет, пустой ли фильтр
func isEmptyFilter(filter *inventory_v1.PartsFilter) bool {
	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

func main() {
	// Создаем gRPC сервер
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Создаем и регистрируем InventoryService
	inventoryService := NewInventoryService()
	inventory_v1.RegisterInventoryServiceServer(grpcServer, inventoryService)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	log.Printf("🚀 InventoryService starting on port %s", grpcPort)
	log.Printf("📦 Available parts: %d", len(inventoryService.parts))

	// Запускаем сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
