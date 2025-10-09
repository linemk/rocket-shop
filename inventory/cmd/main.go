package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
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

	structMetadata, err := structpb.NewStruct(metadata)
	if err != nil {
		log.Printf("Failed to create struct metadata: %v", err)
		return
	}
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

	// Если фильтр пустой, возвращаем все детали
	if req.Filter == nil || isEmptyFilter(req.Filter) {
		result := make([]*inventory_v1.Part, 0, len(s.parts))
		for _, part := range s.parts {
			result = append(result, part)
		}
		return &inventory_v1.ListPartsResponse{Parts: result}, nil
	}

	// Применяем фильтрацию поэтапно
	candidates := s.applyFilters(s.parts, req.Filter)

	// Преобразуем результат в слайс
	result := make([]*inventory_v1.Part, 0, len(candidates))
	for _, part := range candidates {
		result = append(result, part)
	}

	return &inventory_v1.ListPartsResponse{Parts: result}, nil
}

// applyFilters применяет все фильтры к деталям
func (s *InventoryService) applyFilters(parts map[string]*inventory_v1.Part, filter *inventory_v1.PartsFilter) map[string]*inventory_v1.Part {
	candidates := make(map[string]*inventory_v1.Part)

	// Копируем все детали как кандидаты
	for uuid, part := range parts {
		candidates[uuid] = part
	}

	// Применяем фильтры по порядку
	candidates = s.filterByUUIDs(candidates, filter.Uuids)
	candidates = s.filterByNames(candidates, filter.Names)
	candidates = s.filterByCategories(candidates, filter.Categories)
	candidates = s.filterByManufacturerCountries(candidates, filter.ManufacturerCountries)
	candidates = s.filterByTags(candidates, filter.Tags)

	return candidates
}

// filterByUUIDs фильтрует детали по UUID
func (s *InventoryService) filterByUUIDs(candidates map[string]*inventory_v1.Part, uuids []string) map[string]*inventory_v1.Part {
	if len(uuids) == 0 {
		return candidates
	}

	filtered := make(map[string]*inventory_v1.Part)
	for _, uuid := range uuids {
		if part, exists := candidates[uuid]; exists {
			filtered[uuid] = part
		}
	}
	return filtered
}

// filterByNames фильтрует детали по именам
func (s *InventoryService) filterByNames(candidates map[string]*inventory_v1.Part, names []string) map[string]*inventory_v1.Part {
	if len(names) == 0 {
		return candidates
	}

	filtered := make(map[string]*inventory_v1.Part)
	for uuid, part := range candidates {
		for _, name := range names {
			if part.Name == name {
				filtered[uuid] = part
				break
			}
		}
	}
	return filtered
}

// filterByCategories фильтрует детали по категориям
func (s *InventoryService) filterByCategories(candidates map[string]*inventory_v1.Part, categories []inventory_v1.Category) map[string]*inventory_v1.Part {
	if len(categories) == 0 {
		return candidates
	}

	filtered := make(map[string]*inventory_v1.Part)
	for uuid, part := range candidates {
		for _, category := range categories {
			if part.Category == category {
				filtered[uuid] = part
				break
			}
		}
	}
	return filtered
}

// filterByManufacturerCountries фильтрует детали по странам производителей
func (s *InventoryService) filterByManufacturerCountries(candidates map[string]*inventory_v1.Part, countries []string) map[string]*inventory_v1.Part {
	if len(countries) == 0 {
		return candidates
	}

	filtered := make(map[string]*inventory_v1.Part)
	for uuid, part := range candidates {
		if part.Manufacturer == nil {
			continue
		}
		for _, country := range countries {
			if part.Manufacturer.Country == country {
				filtered[uuid] = part
				break
			}
		}
	}
	return filtered
}

// filterByTags фильтрует детали по тегам
func (s *InventoryService) filterByTags(candidates map[string]*inventory_v1.Part, tags []string) map[string]*inventory_v1.Part {
	if len(tags) == 0 {
		return candidates
	}

	filtered := make(map[string]*inventory_v1.Part)
	for uuid, part := range candidates {
		for _, filterTag := range tags {
			for _, partTag := range part.Tags {
				if partTag == filterTag {
					filtered[uuid] = part
					goto nextPart
				}
			}
		}
	nextPart:
	}
	return filtered
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
