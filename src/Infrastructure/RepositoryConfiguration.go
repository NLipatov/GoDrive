package Infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type PersistenceConfig struct {
	EnableAutoDeletion               bool `json:"EnableAutoDeletion"`
	AutoDeletionCheckIntervalMinutes int  `json:"AutoDeletionCheckIntervalMinutes"`
	AutoDeletionIntervalMinutes      int  `json:"AutoDeletionIntervalMinutes"`
}

type AppSettings struct {
	StorageDirectory string            `json:"StorageDirectory"`
	Persistence      PersistenceConfig `json:"Persistence"`
}

func GetConfiguration() AppSettings {
	path, err := os.Getwd()
	if err != nil {
		return defaultConfiguration()
	}

	bytes, err := os.ReadFile(filepath.Join(path, "config.json"))
	if err != nil {
		return defaultConfiguration()
	}

	var appSettings AppSettings
	err = json.Unmarshal(bytes, &appSettings)
	return appSettings
}

func defaultConfiguration() AppSettings {
	return AppSettings{
		StorageDirectory: "storage",
		Persistence: PersistenceConfig{
			EnableAutoDeletion:               true,
			AutoDeletionCheckIntervalMinutes: 30,
			AutoDeletionIntervalMinutes:      180,
		},
	}
}
