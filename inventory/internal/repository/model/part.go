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
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Part struct {
	UUID          string            `bson:"_id"`
	Name          string            `bson:"name"`
	Description   string            `bson:"description"`
	StockQuantity int64             `bson:"stock_quantity"`
	Price         float64           `bson:"price"`
	Category      Category          `bson:"category"`
	Tags          []string          `bson:"tags,omitempty"`
	Metadata      map[string]*Value `bson:"metadata,omitempty"`
	CreatedAt     time.Time         `bson:"created_at"`
	UpdatedAt     *time.Time        `bson:"updated_at,omitempty"`
	Dimensions    *Dimensions       `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer     `bson:"manufacturer,omitempty"`
}

type Value struct {
	String *string  `bson:"string,omitempty"`
	Int64  *int64   `bson:"int64,omitempty"`
	Double *float64 `bson:"double,omitempty"`
	Bool   *bool    `bson:"bool,omitempty"`
}
