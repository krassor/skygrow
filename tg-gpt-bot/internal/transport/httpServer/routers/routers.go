package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/transport/httpServer/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type DeviceRouter struct {
	DeviceHandler handlers.DeviceHandlers
}

func NewDeviceRouter(deviceHandler handlers.DeviceHandlers) *DeviceRouter {
	return &DeviceRouter{
		DeviceHandler: deviceHandler,
	}
}

func (d *DeviceRouter) Router(r *chi.Mux) {
	r.Use(cors.AllowAll().Handler)
	r.Handle("/metrics", promhttp.Handler())
}
