package httpServer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/transport/httpServer/routers"
)

type HttpServer struct {
	Router     *routers.BotRouter
	httpServer *http.Server
}

func NewHttpServer(router *routers.BotRouter) *HttpServer {
	return &HttpServer{
		Router: router,
	}
}

func (h *HttpServer) Listen() {
	app := chi.NewRouter()
	h.Router.Router(app)

	serverPort, ok := os.LookupEnv("HTTP_SERVER_PORT")
	if !ok {
		serverPort = "80"
	}
	serverAddress, ok := os.LookupEnv("HTTP_SERVER_ADDRESS_LISTEN")
	if !ok {
		serverAddress = "0.0.0.0"
	}

	h.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", serverAddress, serverPort),
		Handler: app,
	}

	log.Info().Msgf("Starting http server on %s:%s ...", serverAddress, serverPort)

	err := h.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Warn().Msgf("httpServer.ListenAndServe() Error: %s", err)
	}

	if err == http.ErrServerClosed {
		log.Info().Msgf("%s", err)
	}

}

func (h *HttpServer) Shutdown(ctx context.Context) error {
	if err := h.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
