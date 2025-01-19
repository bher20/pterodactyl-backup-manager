package destinations

import (
	"fmt"
	"os"

	"github.com/bherville/pterodactyl-backup-manager/internal/destinations/azure"
	"github.com/bherville/pterodactyl-backup-manager/internal/destinations/local"
	log "github.com/sirupsen/logrus"
)

type Destinations struct {
	AzureStorageAccount *azure.AzureStorageAccountDestination `json:"azureStorageAccount,omitempty"`
	Local               *local.LocalDestination               `json:"local,omitempty"`
}

func Backup(destinations Destinations, fileName string, fileToBackup *os.File) []error {
	errs := []error{}

	if destinations.AzureStorageAccount != nil {
		log.Debugf("Destination 'Azure Storage Account' enabled, attempting backup...")
		err := azure.Backup(destinations.AzureStorageAccount, fileName, fileToBackup.Name())
		if err != nil {
			log.Error(fmt.Sprintf("error backing up to 'Azure Storage Account' distination: %s", err))
			errs = append(errs, err)
		}
	}

	if destinations.Local != nil {
		log.Debugf("Destination 'Local' enabled, attempting backup...")
		err := local.Backup(destinations.Local, fileName, fileToBackup.Name())
		if err != nil {
			log.Error(fmt.Sprintf("error backing up to 'Local' distination: %s", err))
			errs = append(errs, err)
		}
	}

	return errs
}
