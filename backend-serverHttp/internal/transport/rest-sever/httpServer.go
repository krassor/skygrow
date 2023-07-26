package httpServer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-sever/routers"
)

type HttpServer struct {
	Router     *routers.DeviceRouter
	httpServer *http.Server
}

func NewHttpServer(router *routers.DeviceRouter) *HttpServer {
	return &HttpServer{
		Router: router,
	}
}

func (h *HttpServer) Listen() {
	app := chi.NewRouter()
	h.Router.Router(app)

	serverPort := os.Getenv("DEVICES_HTTP_PORT")
	serverAddress := os.Getenv("DEVICES_HTTP_HOST_LISTEN")
	log.Info().Msgf("Server http get env %s:%s ", serverAddress, serverPort)

	h.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", serverAddress, serverPort),
		Handler: app,
	}
	log.Info().Msgf("Server started on Port %s ", serverPort)

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
