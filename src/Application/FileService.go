package Application

import (
	"GoDrive/src/Domain"
	"github.com/google/uuid"
)

type FileService struct {
	repository FileRepository
}

func NewFileService(repository FileRepository) *FileService {
	return &FileService{repository: repository}
}

func (s *FileService) Get(id uuid.UUID) (*Domain.File, error) {
	file, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (s *FileService) Save(file Domain.File) (uuid.UUID, error) {
	fileId, err := s.repository.Save(file)
	if err != nil {
		return uuid.Nil, err
	}

	return fileId, nil
}

func (s *FileService) Delete(id uuid.UUID) (uuid.UUID, error) {
	fileId, err := s.repository.Delete(id)
	if err != nil {
		return uuid.Nil, err
	}

	return fileId, nil
}
