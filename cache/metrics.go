package cache

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	//cachehit = prometheus.NewCounterVec(
	cachehit = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hit_total",
			Help: "Number of cache hits.",
		},
		[]string{"cache"},
	)
	//cachemiss = prometheus.NewCounterVec(
	cachemiss = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Number of cache miss.",
		},
		[]string{"cache"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	//prometheus.MustRegister(cachehit)
	//prometheus.MustRegister(cachemiss)
}
