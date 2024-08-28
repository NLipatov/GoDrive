package Infrastructure

import (
	"GoDrive/src/Domain"
	"bytes"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileRepository struct{}

func (s *FileRepository) Get(id uuid.UUID) (Domain.File, error) {
	storageDirectoryPath, err := getStorageDirectoryPath()
	if err != nil {
		return Domain.File{}, err
	}

	file, err := os.Open(filepath.Join(storageDirectoryPath, id.String()))
	if err != nil {
		return Domain.File{}, err
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return Domain.File{}, err
	}

	return Domain.File{
		Id:        id,
		Data:      fileContent,
		ExpiresAt: time.Time{},
	}, nil
}

func (s *FileRepository) Save(file Domain.File) (uuid.UUID, error) {
	storageDirectoryPath, err := getStorageDirectoryPath()
	if err != nil {
		return uuid.Nil, err
	}

	for {
		file.Id = uuid.New()

		filePath := filepath.Join(storageDirectoryPath, file.Id.String())

		if _, err = os.Stat(filePath); os.IsNotExist(err) {
			break
		}
	}

	writeToFile, err := os.Create(filepath.Join(storageDirectoryPath, file.Id.String()))
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
	storageDirectoryPath, err := getStorageDirectoryPath()
	if err != nil {
		return uuid.Nil, err
	}

	filePath := filepath.Join(storageDirectoryPath, id.String())
	err = os.Remove(filePath)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Nil, nil
}

func getStorageDirectoryPath() (string, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	storageDirectory := filepath.Join(currentDirectory, "storage")
	if _, err = os.Stat(storageDirectory); os.IsNotExist(err) {
		err = os.Mkdir(storageDirectory, 0755)
		if err != nil {
			return "", err
		}
	}

	return storageDirectory, nil
}
