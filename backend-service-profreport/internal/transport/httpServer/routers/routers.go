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
			mux.Route("/questionnaire", func(mux chi.Router) {
				//Private
				mux.Group(func(mux chi.Router) {
					mux.Post("/adult", r.questionnaireHandler.AdultCreate)
				},
				)
				mux.Group(func(mux chi.Router) {
					mux.Post("/schoolchild", r.questionnaireHandler.SchoolchildCreate)
				},
				)
			},
			)

			mux.Route("/promocode", func(mux chi.Router) {
				mux.Post("/apply", r.questionnaireHandler.ApplyPromoCode)
			},
			)

			mux.Route("/prices", func(mux chi.Router) {
				mux.Get("/", r.questionnaireHandler.GetTestPrices)
			},
			)

			mux.Route("/callback", func(mux chi.Router) {
				mux.Route("/cloudpayments", func(mux chi.Router) {
					mux.Group(func(mux chi.Router) {
						mux.Post("/pay", r.questionnaireHandler.Payment)
					},
					)
				},
				)
			},
			)
		},
		)
	},
	)

}
