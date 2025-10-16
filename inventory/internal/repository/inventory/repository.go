package inventory

import (
	"context"
	"fmt"
	"sync"

	"github.com/linemk/rocket-shop/inventory/internal/entyties/models"
	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

type Repository struct {
	mu    sync.RWMutex
	parts map[string]models.Part
}

func NewRepository() *Repository {
	return &Repository{
		parts: make(map[string]models.Part),
	}
}

func (r *Repository) GetPart(ctx context.Context, uuid string) (models.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid]
	if !exists {
		return models.Part{}, fmt.Errorf("part with UUID %s not found", uuid)
	}

	return part, nil
}

func (r *Repository) ListParts(ctx context.Context, filter models.PartFilter) ([]models.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Если фильтр пустой, возвращаем все детали
	if isEmptyFilter(filter) {
		result := make([]models.Part, 0, len(r.parts))
		for _, part := range r.parts {
			result = append(result, part)
		}
		return result, nil
	}

	// Применяем фильтрацию
	candidates := r.applyFilters(r.parts, filter)

	// Преобразуем результат в слайс
	result := make([]models.Part, 0, len(candidates))
	for _, part := range candidates {
		result = append(result, part)
	}

	return result, nil
}

func (r *Repository) CreatePart(ctx context.Context, part models.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.parts[part.UUID]; exists {
		return fmt.Errorf("part with UUID %s already exists", part.UUID)
	}

	r.parts[part.UUID] = part
	return nil
}

func (r *Repository) UpdatePart(ctx context.Context, uuid string, part models.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.parts[uuid]; !exists {
		return fmt.Errorf("part with UUID %s not found", uuid)
	}

	r.parts[uuid] = part
	return nil
}

func (r *Repository) DeletePart(ctx context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.parts[uuid]; !exists {
		return fmt.Errorf("part with UUID %s not found", uuid)
	}

	delete(r.parts, uuid)
	return nil
}

// applyFilters применяет все фильтры к деталям
func (r *Repository) applyFilters(parts map[string]models.Part, filter models.PartFilter) map[string]models.Part {
	candidates := make(map[string]models.Part)

	// Копируем все детали как кандидаты
	for uuid, part := range parts {
		candidates[uuid] = part
	}

	// Применяем фильтры по порядку
	candidates = r.filterByUUIDs(candidates, filter.UUIDs)
	candidates = r.filterByNames(candidates, filter.Names)
	candidates = r.filterByCategories(candidates, filter.Categories)
	candidates = r.filterByManufacturerCountries(candidates, filter.ManufacturerCountries)
	candidates = r.filterByTags(candidates, filter.Tags)

	return candidates
}

// filterByUUIDs фильтрует детали по UUID
func (r *Repository) filterByUUIDs(candidates map[string]models.Part, uuids []string) map[string]models.Part {
	if len(uuids) == 0 {
		return candidates
	}

	filtered := make(map[string]models.Part)
	for _, uuid := range uuids {
		if part, exists := candidates[uuid]; exists {
			filtered[uuid] = part
		}
	}
	return filtered
}

// filterByNames фильтрует детали по именам
func (r *Repository) filterByNames(candidates map[string]models.Part, names []string) map[string]models.Part {
	if len(names) == 0 {
		return candidates
	}

	filtered := make(map[string]models.Part)
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
func (r *Repository) filterByCategories(candidates map[string]models.Part, categories []inventory_v1.Category) map[string]models.Part {
	if len(categories) == 0 {
		return candidates
	}

	filtered := make(map[string]models.Part)
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
func (r *Repository) filterByManufacturerCountries(candidates map[string]models.Part, countries []string) map[string]models.Part {
	if len(countries) == 0 {
		return candidates
	}

	filtered := make(map[string]models.Part)
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
func (r *Repository) filterByTags(candidates map[string]models.Part, tags []string) map[string]models.Part {
	if len(tags) == 0 {
		return candidates
	}

	filtered := make(map[string]models.Part)
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
func isEmptyFilter(filter models.PartFilter) bool {
	return len(filter.UUIDs) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}
