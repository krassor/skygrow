package httpServer

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/krassor/skygrow/backend-service-auth/internal/config"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/krassor/skygrow/backend-service-auth/internal/transport/rest-server/routers"
)

type HttpServer struct {
	Router     *routers.Router
	httpServer *http.Server
}

func NewHttpServer(router *routers.Router) *HttpServer {
	return &HttpServer{
		Router: router,
	}
}

func (h *HttpServer) Listen(cfg *config.Config) {
	app := chi.NewRouter()
	h.Router.Router(app)

	serverPort := cfg.HttpServer.Port
	serverAddress := cfg.HttpServer.Address
	log.Info().Msgf("Server http get env %s:%s ", serverAddress, serverPort)

	h.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", serverAddress, serverPort),
		Handler: app,
	}
	log.Info().Msgf("Server started on Port %s ", serverPort)

	err := h.httpServer.ListenAndServe()

	if err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Warn().Msgf("httpServer.ListenAndServe() Error: %s", err)
	}

	if errors.Is(err, http.ErrServerClosed) {
		log.Info().Msgf("%s", err)
	}

}

func (h *HttpServer) Shutdown(ctx context.Context) error {
	if err := h.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
