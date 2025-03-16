package metric

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	RequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "total number of request for api with status codes",
		},
		[]string{"method", "status"},
	)
	CacheSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "total numbers of elems in cache",
		},
		[]string{"cachename"},
	)
)

var once sync.Once

func MustRegisterMetrics() {
	once.Do(func() {
		prometheus.MustRegister(RequestsCounter, CacheSize)
	})
}

type CacheStatGetter interface {
	GetBlogLen() (string, int)
	GetPostLen() (string, int)
}

func GoCountCacheLen(ctx context.Context, interval time.Duration, cacheStats CacheStatGetter) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
				name, len := cacheStats.GetBlogLen()
				log.Info().Str("blogcache", name).Int("len", len).Msg("")
				CacheSize.WithLabelValues(name).Set(float64(len))
				name, len = cacheStats.GetPostLen()
				log.Info().Str("postcache", name).Int("len", len).Msg("")
				CacheSize.WithLabelValues(name).Set(float64(len))
			}
		}
	}()
}
