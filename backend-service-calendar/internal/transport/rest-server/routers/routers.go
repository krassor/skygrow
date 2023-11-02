package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"

	//"github.com/go-chi/cors"
	"github.com/krassor/skygrow/backend-service-calendar/internal/transport/rest-server/handlers"
	myMiddleware "github.com/krassor/skygrow/backend-service-calendar/internal/transport/rest-server/middleware"
)

type Router struct {
	BookOrderHandler *handlers.BookOrderHandler
	UserHandler      *handlers.UserHandler
	tokenAuth        *jwtauth.JWTAuth
}

func NewRouter(bookOrderHandler *handlers.BookOrderHandler, userHandler *handlers.UserHandler) *Router {
	return &Router{
		BookOrderHandler: bookOrderHandler,
		UserHandler:      userHandler,
		tokenAuth:        jwtauth.New("HS256", []byte("skygrowSecretKey"), nil),
	}
}

func (d *Router) Router(r *chi.Mux) {
	r.Use(cors.AllowAll().Handler)
	r.Use(myMiddleware.LoggerMiddleware)
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/user", func(r chi.Router) {
		//Public
		r.Group(func(r chi.Router) {
			r.Post("/signup", d.UserHandler.SignUp)
			r.Post("/signin", d.UserHandler.SignIn)
		})

		//Private
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(d.tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Post("/bookorder", d.BookOrderHandler.CreateBookOrder)
		})
	})

}
