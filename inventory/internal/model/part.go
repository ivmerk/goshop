package model

import (
	"time"
)

type Part struct {
	UUID          string            `bson:"_id"`
	Name          string            `bson:"name"`
	Description   string            `bson:"description"`
	StockQuantity int64             `bson:"stock_quantity"`
	Price         float64           `bson:"price"`
	Category      Category          `bson:"category"`
	Tags          []string          `bson:"tags"`
	Metadata      map[string]*Value `bson:"metadata"`
	CreatedAt     time.Time         `bson:"created_at"`
	UpdatedAt     *time.Time        `bson:"updated_at,omitempty"`
	Dimensions    *Dimensions       `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer     `bson:"manufacturer,omitempty"`
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Category string

const (
	CategoryUnknown  Category = "UNKNOWN"
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
)

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Value struct {
	String *string
	Int64  *int64
	Double *float64
	Bool   *bool
}
