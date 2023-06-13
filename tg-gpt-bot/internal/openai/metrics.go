package openai

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var TotalTokensUsage = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace: "tgGptBot",
	Subsystem: "openai",
	Name:      "total_token_usage_per_user",
	Help:      "Количество токенов в запросе пользователя",
}, []string{"username"})

func observeTotalTokensUsage(totalTokenUsage int, username string) {
	TotalTokensUsage.WithLabelValues(username).Observe(float64(totalTokenUsage))
}
