package main

import (
	"github.com/krassor/skygrow/backend-service-calendar/internal/services/GoogleService"
	"github.com/rs/zerolog/log"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(2)
	gc1 := GoogleService.NewGoogleCalendar()
	gc2 := GoogleService.NewGoogleCalendar()
	log.Info().Msg("GoogleService created")
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		id, err := gc1.CreateCalendar(
			"New course for learning",
			"Smetankin Dmirtrii",
			"Europe/Moscow")

		if err != nil {
			log.Error().Msgf("ERROR: %v", err)
		}
		log.Info().Msgf("calendarID: %s", id)
	}()

	go func() {
		defer wg.Done()
		id, err := gc2.CreateCalendar(
			"New course for learning",
			"Smetankina Lubov",
			"Europe/Moscow")

		if err != nil {
			log.Error().Msgf("ERROR: %v", err)
		}
		log.Info().Msgf("calendarID: %s", id)
	}()
	wg.Wait()
}
