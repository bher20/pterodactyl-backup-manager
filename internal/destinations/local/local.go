package local

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type LocalDestination struct {
	BackupDirectory string `json:"backupDirectory"`
}

func Backup(local *LocalDestination, fileName string, filePath string) error {
	_, err := os.Stat(local.BackupDirectory)
	if err != nil {
		return err
	}

	backupPath := filepath.Join(local.BackupDirectory, fileName)
	log.Trace(fmt.Sprintf("Writing backup to: %s", backupPath))

	err = os.Rename(filePath, backupPath)
	if err != nil {
		return err
	}

	return nil
}
