package part

import (
	"time"

	"github.com/google/uuid"

	repoModel "github.com/ivmerk/goshop/inventory/internal/repository/model"
)

func (r *repository) seed() {
	now := time.Now()

	add := func(p *repoModel.Part) {
		p.UUID = uuid.NewString()
		p.CreatedAt = now
		r.parts[p.UUID] = p
	}

	add(&repoModel.Part{
		Name:          "Main booster",
		Description:   "Основной маршевый двигатель",
		Price:         12500.50,
		StockQuantity: 5,
		Category:      repoModel.CategoryEngine,
		Dimensions:    &repoModel.Dimensions{Length: 300, Width: 120, Height: 120, Weight: 850},
		Manufacturer:  &repoModel.Manufacturer{Name: "Rocket GmbH", Country: "Germany", Website: "https://rocket.example"},
		Tags:          []string{"engine", "main"},
		Metadata: map[string]*repoModel.Value{
			"series": {String: ptr("MB-1")},
			"thrust": {Double: ptr(1200.5)},
			"tested": {Bool: ptr(true)},
			"batch":  {Int64: ptr(int64(42))},
		},
	})
	add(&repoModel.Part{
		Name:          "Fuel tank",
		Description:   "Топливный бак на 5000 литров",
		Price:         3200,
		StockQuantity: 12,
		Category:      repoModel.CategoryFuel,
		Dimensions:    &repoModel.Dimensions{Length: 200, Width: 150, Height: 150, Weight: 400},
		Manufacturer:  &repoModel.Manufacturer{Name: "Tank Corp", Country: "USA", Website: "https://tank.example"},
		Tags:          []string{"fuel", "tank"},
	})
	add(&repoModel.Part{
		Name:          "Panoramic porthole",
		Description:   "Панорамный иллюминатор",
		Price:         980.99,
		StockQuantity: 30,
		Category:      repoModel.CategoryPorthole,
		Dimensions:    &repoModel.Dimensions{Length: 80, Width: 80, Height: 10, Weight: 25},
		Manufacturer:  &repoModel.Manufacturer{Name: "Glass Ltd", Country: "Japan", Website: "https://glass.example"},
		Tags:          []string{"porthole", "glass"},
	})
	add(&repoModel.Part{
		Name:          "Delta wing",
		Description:   "Дельтовидное крыло",
		Price:         7400,
		StockQuantity: 8,
		Category:      repoModel.CategoryWing,
		Dimensions:    &repoModel.Dimensions{Length: 600, Width: 250, Height: 30, Weight: 300},
		Manufacturer:  &repoModel.Manufacturer{Name: "Rocket GmbH", Country: "Germany", Website: "https://rocket.example"},
		Tags:          []string{"wing", "aero"},
	})
}

func ptr[T any](v T) *T { return &v }
