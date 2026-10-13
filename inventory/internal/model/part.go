package model

import "time"

type Part struct {
	UUID          string
	Name          string
	Description   string
	StockQuantity int64
	Price         float64
	Category      Category
	Tags          []string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
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
