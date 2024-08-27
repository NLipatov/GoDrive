package Application

import (
	"GoDrive/src/Domain"
	"github.com/google/uuid"
)

type FileRepository interface {
	Get(id uuid.UUID) (Domain.File, error)
	Save(file Domain.File) (uuid.UUID, error)
	Delete(id uuid.UUID) (uuid.UUID, error)
}
