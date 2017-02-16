package lazycache

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	PromCacheRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_request_total",
			Help: "Number of cache requests.",
		},
    []string{"store"},
	)

	PromCacheMisses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Number of cache misses.",
		},
    []string{"store"},
	)

	PromCacheSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cache_size",
		Help: "",
	},
  []string{"store"},
)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(PromCacheRequests)
	prometheus.MustRegister(PromCacheMisses)
	prometheus.MustRegister(PromCacheSize)
}
