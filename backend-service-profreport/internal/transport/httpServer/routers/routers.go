package routers

import (
	"app/main.go/internal/transport/httpServer/handlers"
	myMiddleware "app/main.go/internal/transport/httpServer/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	questionnaireHandler *handlers.QuestionnaireHandler
	secret               string
}

func NewRouter(questionnaireHandler *handlers.QuestionnaireHandler, secret string) *Router {
	return &Router{
		questionnaireHandler: questionnaireHandler,
		secret:               secret,
	}
}

func (r *Router) Mount(mux *chi.Mux) {

	mux.Use(cors.AllowAll().Handler)
	mux.Use(myMiddleware.LoggerMiddleware)
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/v1", func(mux chi.Router) {

			//Private
			mux.Group(func(mux chi.Router) {
				//Для отладки временно без авторизации
				//mux.Use(myMiddleware.Authorization(r.secret))
				mux.Post("/questionnaire", r.questionnaireHandler.Create)
			})
		},
		)
	},
	)

}
