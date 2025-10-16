package models

import (
	"time"

	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь в инвентаре
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      inventory_v1.Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Dimensions представляет размеры детали
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

// Manufacturer представляет производителя
type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// PartFilter представляет фильтр для поиска деталей
type PartFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []inventory_v1.Category
	ManufacturerCountries []string
	Tags                  []string
}

// PartInfo представляет информацию о детали для клиентов
type PartInfo struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      inventory_v1.Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
