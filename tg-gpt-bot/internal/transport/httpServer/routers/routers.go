package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	//"github.com/krassor/skygrow/tg-gpt-bot/internal/transport/httpServer/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type BotRouter struct {
	//DeviceHandler handlers.DeviceHandlers
}

func NewBotRouter( /*deviceHandler handlers.DeviceHandlers*/ ) *BotRouter {
	return &BotRouter{
		//DeviceHandler: deviceHandler,
	}
}

func (br *BotRouter) Router(r *chi.Mux) {
	r.Use(cors.AllowAll().Handler)
	r.Handle("/metrics", promhttp.Handler())
}
