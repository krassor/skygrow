package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/cors"
	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-server/handlers"
	
)

type Router struct {
	BookOrderHandler *handlers.BookOrderHandler
	UserHandler *handlers.UserHandler
}

func NewRouter(bookOrderHandler *handlers.BookOrderHandler, userHandler *handlers.UserHandler) *Router {
	return &Router{
		BookOrderHandler: bookOrderHandler,
		UserHandler: userHandler,
	}
}

func (d *Router) Router(r *chi.Mux) {
	//r.Use(cors.AllowAll().Handler)
	r.Use(middleware.Heartbeat("/ping"))
	r.Post("/bookorder", d.BookOrderHandler.CreateBookOrder)
	r.Post("/signup", d.UserHandler.SignUp)
}
