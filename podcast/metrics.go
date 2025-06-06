package podcast

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	rssGetSecondsMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "rss_get_seconds",
		Help:    "Time to get an RSS feed",
		Buckets: []float64{0.1, 1, 5, 10, 20, 30},
	}, []string{"cached"})
	rssGetTotalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rss_get_total",
		Help: "Number of fetched RSS feeds",
	}, []string{"program_id"})
	hashLookup = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "hash_lookup",
		Help: "Number of times a hash was provided to and used to lookup cached podcasts",
	}, []string{"success"})
)
