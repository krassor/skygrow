package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/dto"
	services "github.com/krassor/skygrow/backend-serverHttp/internal/services/devices"
	"github.com/krassor/skygrow/backend-serverHttp/pkg/utils"
)

type DeviceHandlers interface {
	CreateDevice(w http.ResponseWriter, r *http.Request)
}

type deviceHandler struct {
	deviceService services.DevicesRepoService
}

func NewDeviceHandler(deviceService services.DevicesRepoService) DeviceHandlers {
	return &deviceHandler{
		deviceService: deviceService,
	}
}

func (d *deviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	deviceDto := dto.RequestDeviceDto{}
	err := json.NewDecoder(r.Body).Decode(&deviceDto)
	if err != nil {
		log.Warn().Msgf("Error decode json: %s", err)
		httpErr := utils.Err(w, http.StatusInternalServerError, err)
		if httpErr != nil {
			log.Warn().Msgf("Cannot sending error message to http client: %s", httpErr)
		}
		return
	}
	deviceUrl := fmt.Sprintf("%s://%s:%s", deviceDto.DeviceSchema, deviceDto.DeviceIpAddress, deviceDto.DevicePort)
	_, err = url.Parse(deviceUrl)
	if err != nil {
		log.Warn().Msgf("Error parse URL: %s", err)
		httpErr := utils.Err(w, http.StatusInternalServerError, err)
		if httpErr != nil {
			log.Warn().Msgf("Cannot sending error message to http client: %s", httpErr)
		}
		return
	}
	//deviceEntity := entities.Devices{}
	deviceEntity, err := d.deviceService.CreateNewDevice(r.Context(), deviceDto)

	if err != nil {
		log.Error().Msgf("Error creating device: %s", err)
		httpErr := utils.Err(w, http.StatusInternalServerError, err)
		if httpErr != nil {
			log.Warn().Msgf("Cannot sending error message to http client: %s", httpErr)
		}
		return
	}

	responseDeviceParams := dto.ResponseDeviceParams{}
	responseDeviceParams.DeviceId = deviceEntity.ID
	responseDevice := utils.Message(true, responseDeviceParams)
	err = utils.Json(w, http.StatusOK, responseDevice)
	if err != nil {
		log.Warn().Msgf("Cannot encode response to json: %s", err)
	}
}
