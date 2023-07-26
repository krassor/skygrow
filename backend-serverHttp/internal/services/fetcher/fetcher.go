package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
	services "github.com/krassor/skygrow/backend-serverHttp/internal/services/devices"
	telegramBot "github.com/krassor/skygrow/backend-serverHttp/internal/telegram"
	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-client/client"
)

const (
	fetcherDuration      time.Duration = 35 * time.Second
	timeoutGetDeviceList time.Duration = 3 * time.Second
)

type DeviceFetcher struct {
	client  client.DeviceStatusClient
	service services.DevicesRepoService
	bot     *telegramBot.Bot
}

func NewDeviceFetcher(service services.DevicesRepoService, bot *telegramBot.Bot) *DeviceFetcher {
	return &DeviceFetcher{client: client.NewDefaultDevice(&http.Client{}), service: service, bot: bot}
}

func (f *DeviceFetcher) Start(ctx context.Context) {

	var wg sync.WaitGroup
	for {
		entityList := f.getDeviceList(ctx, timeoutGetDeviceList)
		log.Info().Msgf("select default")

		wg.Add(len(entityList))
		log.Info().Msgf("workgroup add %d elements", len(entityList))

		for _, e := range entityList {
			go f.requestDeviceStatus(ctx, e, &wg)
		}
		wg.Wait()
		time.Sleep(fetcherDuration)

	}
}

func (f *DeviceFetcher) requestDeviceStatus(ctx context.Context, entity entities.Devices, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Msgf("Enter anonymous go routine")

	deviceUrl := fmt.Sprintf("%s://%s:%s", entity.DeviceSchema, entity.DeviceIpAddress, entity.DevicePort)
	log.Info().Msgf("device URL: %s", deviceUrl)

	status, err := f.client.GetStatus(ctx, deviceUrl)
	if err != nil {
		log.Error().Msgf("Error get device %d %s %s status in fetcher.Start(): %s", entity.ID, entity.DeviceVendor, entity.DeviceName, err)
	}
	log.Info().Msgf("Device status: %t", status)

	if status != entity.DeviceStatus {
		err = f.bot.DeviceStatusNotify(ctx, entity, status)
		if err != nil {
			log.Error().Msgf("Error notify subscribers: %s", err)
		}
		_, err = f.service.UpdateDeviceStatus(ctx, entity, status)
		if err != nil {
			log.Error().Msgf("Error update device %d %s %s status in fetcher.Start(): %s", entity.ID, entity.DeviceVendor, entity.DeviceName, err)
		}

		log.Info().Msgf("End of anonymous go routine")
	}

}

func (f *DeviceFetcher) getDeviceList(ctx context.Context, timeoutDuration time.Duration) []entities.Devices {
	var entityList []entities.Devices
	var err error
	for {
		entityList, err = f.service.GetDevices(ctx)
		if err != nil {
			log.Error().Msgf("Fetcher: cannot get data from repo: %s", err)
			time.Sleep(timeoutDuration)
			continue
		}

		break
	}
	return entityList
}
