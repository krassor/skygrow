package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-server/handlers"
	
)

type Router struct {
	Handler *handlers.BookOrderHandler
}

func NewRouter(handler *handlers.BookOrderHandler) *Router {
	return &Router{
		Handler: handler,
	}
}

func (d *Router) Router(r *chi.Mux) {
	r.Use(cors.AllowAll().Handler)
	r.Post("/bookorder", d.Handler.CreateBookOrder)
}
