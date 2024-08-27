package Infrastructure

import (
	"GoDrive/src/Domain"
	"bytes"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

type FileRepository struct{}

func (s *FileRepository) Get(id uuid.UUID) (Domain.File, error) {
	panic("implement me")
}

func (s *FileRepository) Save(file Domain.File) (uuid.UUID, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return uuid.Nil, err
	}

	storageDirectory := filepath.Join(currentDirectory, "storage")
	if _, err = os.Stat(storageDirectory); os.IsNotExist(err) {
		err = os.Mkdir(storageDirectory, 0755)
		if err != nil {
			return uuid.Nil, err
		}
	}

	for {
		file.Id = uuid.New()

		filePath := filepath.Join(storageDirectory, file.Id.String())

		if _, err = os.Stat(filePath); os.IsNotExist(err) {
			break
		}
	}

	writeToFile, err := os.Create(filepath.Join(storageDirectory, file.Id.String()))
	if err != nil {
		return uuid.Nil, err
	}
	defer writeToFile.Close()

	_, err = io.Copy(writeToFile, bytes.NewReader(file.Data))
	if err != nil {
		return uuid.Nil, err
	}

	return file.Id, nil
}

func (s *FileRepository) Delete(id uuid.UUID) (uuid.UUID, error) {
	panic("implement me")
}
