package valueobject

import (
	"github.com/google/uuid"
	"time"
)

type Measuring struct {
	id uuid.UUID
	fromDevice uuid.UUID
	timeStamp time.Time
}