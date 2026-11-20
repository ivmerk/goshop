// Package metrics — бизнес-метрики Assembly-сервиса.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// AssemblyDuration — длительность сборки каждого корабля.
var AssemblyDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "assembly_duration_seconds",
	Help:    "Длительность сборки корабля.",
	Buckets: []float64{1, 2, 3, 4, 5, 6, 8, 10, 15, 30},
})
