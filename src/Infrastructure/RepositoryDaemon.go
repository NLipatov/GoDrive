package Infrastructure

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type RepositoryDaemon struct {
	AppSettings AppSettings
}

func NewRepositoryDaemon() *RepositoryDaemon {
	return &RepositoryDaemon{}
}

func (d *RepositoryDaemon) WatchStorageDirectory() {
	go func() {
		for {
			config := GetConfiguration()

			d.AppSettings = config

			if !config.Persistence.EnableAutoDeletion {
				time.Sleep(time.Duration(config.Persistence.AutoDeletionCheckIntervalMinutes) * time.Minute)
				continue
			}

			err := d.cleanupDirectory(config.StorageDirectory)
			if err != nil {
				log.Printf("Daemon failed to cleanup storage: %v", err.Error())
			}

			time.Sleep(time.Duration(config.Persistence.AutoDeletionCheckIntervalMinutes) * time.Minute)
		}
	}()
}

func (d *RepositoryDaemon) cleanupDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			err = d.cleanupDirectory(path)
			if err != nil {
				log.Printf("Error while attempting to cleanup directory: %v", err)
				continue
			}

			if isEmptyDir, err := isDirEmpty(path); err == nil && isEmptyDir {
				err := os.Remove(path)
				if err != nil {
					log.Printf("Failed to delete directory: %v", err)
				}
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				log.Printf("Failed to get file info: %v", err)
				continue
			}

			if time.Since(info.ModTime()) > time.Duration(d.AppSettings.Persistence.AutoDeletionIntervalMinutes)*time.Minute {
				err = os.Remove(path)
				if err != nil {
					log.Printf("Failed to delete file: %v", err)
				}
			}
		}
	}

	return nil
}

func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
