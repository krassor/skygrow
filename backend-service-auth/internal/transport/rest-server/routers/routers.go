package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"

	//"github.com/go-chi/cors"
	"github.com/krassor/skygrow/backend-service-auth/internal/transport/rest-server/handlers"
)

type Router struct {
	UserHandler *handlers.UserHandler
	tokenAuth   *jwtauth.JWTAuth
}

func NewRouter(userHandler *handlers.UserHandler) *Router {
	return &Router{
		UserHandler: userHandler,
		tokenAuth:   jwtauth.New("HS256", []byte("skygrowSecretKey"), nil),
	}
}

func (d *Router) Router(r *chi.Mux) {
	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/user", func(r chi.Router) {
		//Public
		r.Group(func(r chi.Router) {
			r.Post("/signup", d.UserHandler.SignUp)
			r.Post("/signin", d.UserHandler.SignIn)
		})

		//Private
		//r.Group(func(r chi.Router) {
		//	r.Use(jwtauth.Verifier(d.tokenAuth))
		//	r.Use(jwtauth.Authenticator)
		//	r.Post("/<privateEndpoint>", /*privateHandler*/)
		//})
	})

}
