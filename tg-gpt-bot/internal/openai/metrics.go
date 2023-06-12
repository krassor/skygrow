package openai

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "tg-gpt-bot",
	Subsystem:  "openai",
	Name:       "requestTime",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"username"})

func observeRequest(d time.Duration, username string) {
	requestMetrics.WithLabelValues(username).Observe(d.Seconds())
}