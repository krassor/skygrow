package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	retryCount    int           = 3
	retryDuration time.Duration = 3 * time.Second
)

type DefaultDevice struct {
	httpClient *http.Client
}

func NewDefaultDevice(httpClient *http.Client) DeviceStatusClient {
	return &DefaultDevice{httpClient: httpClient}
}

func (d *DefaultDevice) GetStatus(ctx context.Context, url string) (bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request with ctx: %w", err)
	}
	//d.httpClient.Timeout = 500 * time.Millisecond

	var res *http.Response
	for i := 0; i < retryCount; i++ {
		r, e := d.httpClient.Do(req)
		if e != nil {
			err = e
			log.Error().Msgf("Error try %s:\n%w", i+1, err)
			time.Sleep(retryDuration)
			continue
		}
		res = r
		err = e
		break
	}
	if err != nil {
		return false, fmt.Errorf("failed to perform http request: %w", err)
	}
	//if res.StatusCode != http.StatusOK {
	//	return false, fmt.Errorf("status code is not 200 OK")
	//}
	if res.ContentLength == 0 {
		return false, fmt.Errorf("response body is null")
	}

	return true, nil
}
