package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	//"github.com/krassor/skygrow/tg-gpt-bot/internal/transport/httpServer/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type BotRouter struct {
	//openAiHandler handlers.openAiHandler
}

func NewBotRouter( /*deviceHandler handlers.DeviceHandlers*/ ) *BotRouter {
	return &BotRouter{
		//DeviceHandler: deviceHandler,
	}
}

func (br *BotRouter) Router(r *chi.Mux) {

	r.Use(middleware.Heartbeat("/ping"))

	//Public
	r.Group(func(r chi.Router) {
		r.Handle("/metrics", promhttp.Handler())
	})

	// r.Route("/api", func(r chi.Router) {
	// 	//Private
	// 	r.Group(func(r chi.Router) {
	// 		r.Use(middleware.Logger)
	// 		//r.Use(middleware.BasicAuth())
	// 		r.Post("/createchatcompletion", br.BookOrderHandler.CreateBookOrder)
	// 	})
	// })

}
