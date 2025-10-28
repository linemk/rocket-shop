package models

import (
	"time"

	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь в инвентаре
type Part struct {
	UUID          string                 `bson:"uuid"`
	Name          string                 `bson:"name"`
	Description   string                 `bson:"description"`
	Price         float64                `bson:"price"`
	StockQuantity int64                  `bson:"stock_quantity"`
	Category      inventory_v1.Category  `bson:"category"`
	Dimensions    *Dimensions            `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer          `bson:"manufacturer,omitempty"`
	Tags          []string               `bson:"tags,omitempty"`
	Metadata      map[string]interface{} `bson:"metadata,omitempty"`
	CreatedAt     time.Time              `bson:"created_at"`
	UpdatedAt     time.Time              `bson:"updated_at"`
}

// Dimensions представляет размеры детали
type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

// Manufacturer представляет производителя
type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
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
