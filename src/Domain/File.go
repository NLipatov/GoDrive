package Domain

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	Id        uuid.UUID
	Data      []byte
	ExpiresAt time.Time
}
