package aggregate

import (
	"github.com/krassor/skygrow/internal/entity"
	"github.com/krassor/skygrow/internal/valueobject"

)

type Growbox struct {
	account *entity.Account
	devices []*entity.Device
	measuring []*valueobject.Measuring
}