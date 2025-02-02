package telegramBot

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ResponseLatencySecSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "tgGptBot",
	Subsystem:  "telegram",
	Name:       "request_latency_sec",
	Help:       "Время задержки обработки запроса в секундах",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"username"})

func observeResponseLatencySecSummary(d time.Duration, username string) {
	ResponseLatencySecSummary.WithLabelValues(username).Observe(d.Seconds())
}
