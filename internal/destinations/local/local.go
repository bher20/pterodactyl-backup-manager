package local

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type LocalDestination struct {
	BackupDirectory string       `json:"backupDirectory"`
	FilePermissions *os.FileMode `json:"filePermissions,omitempty"`
}

func Backup(local *LocalDestination, fileName string, fileData []byte) error {
	permissions := os.FileMode(0644)
	if local.FilePermissions != nil {
		permissions = *local.FilePermissions
	}

	_, err := os.Stat(local.BackupDirectory)
	if err != nil {
		return err
	}

	backupPath := filepath.Join(local.BackupDirectory, fileName)
	log.Trace(fmt.Sprintf("Writing backup to: %s", backupPath))

	return os.WriteFile(backupPath, fileData, permissions)
}
