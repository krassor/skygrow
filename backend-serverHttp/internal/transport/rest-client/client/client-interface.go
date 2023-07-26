package client

import (
	"context"
)

type DeviceStatusClient interface {
	GetStatus(ctx context.Context, url string) (bool, error)
}
