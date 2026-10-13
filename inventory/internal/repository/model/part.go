package model

import "time"

type Category string

const (
	CategoryUnknown  Category = "UNKNOWN"
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
)

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

type Part struct {
	UUID          string
	Name          string
	Description   string
	StockQuantity int64
	Price         float64
	Category      Category
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
}

type Value struct {
	String *string
	Int64  *int64
	Double *float64
	Bool   *bool
}
