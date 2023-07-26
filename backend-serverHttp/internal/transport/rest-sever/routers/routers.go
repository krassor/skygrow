package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-sever/handlers"
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
	r.Post("/devices", d.DeviceHandler.CreateDevice)
}
