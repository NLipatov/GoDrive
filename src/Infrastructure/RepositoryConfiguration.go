package Infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type PersistenceConfig struct {
	EnableAutoDeletion          bool `json:"EnableAutoDeletion"`
	AutoDeletionIntervalMinutes int  `json:"AutoDeletionIntervalMinutes"`
}

type AppSettings struct {
	StorageDirectory string            `json:"StorageDirectory"`
	Persistence      PersistenceConfig `json:"Persistence"`
}

func GetConfiguration() (AppSettings, error) {
	path, err := os.Getwd()
	if err != nil {
		return AppSettings{}, err
	}

	bytes, err := os.ReadFile(filepath.Join(path, "config.json"))
	if err != nil {
		return AppSettings{}, err
	}

	var appSettings AppSettings
	err = json.Unmarshal(bytes, &appSettings)
	return appSettings, nil
}
