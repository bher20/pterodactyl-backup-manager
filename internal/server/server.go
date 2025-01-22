package server

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bherville/pterodactyl-backup-manager/internal/config"
	"github.com/bherville/pterodactyl-backup-manager/internal/destinations"
	"github.com/bherville/pterodactyl-backup-manager/internal/utils"

	"github.com/bherville/pterodactyl-sdk-go/pkg/pterodactyl"

	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	cronSchedular *cron.Cron
)

const (
	WAIT_FOR_BACKUP_SECONDS = 2
)

func PerformBackup(backupConfig *config.Backup, pterodactylServer *pterodactyl.PterodactylServer, appServer pterodactyl.Server, errorsFatal bool, tmpDirPath string, deleteAfterBackup bool) (*pterodactyl.Backup, error) {
	log.Info(fmt.Sprintf("Running backup of '%s'...", appServer.Attributes.Name))

	backup, err := pterodactyl.BackupServerWithWait(*pterodactylServer, appServer)
	if err != nil {
		return nil, err
	}

	downloadPath := filepath.Join(tmpDirPath, backup.Attributes.UUID)
	log.Info(fmt.Sprintf("Downlading new backup '%s'!", backup.Attributes.UUID))
	backupFile, err := pterodactyl.DownloadServerBackup(*pterodactylServer, appServer, backup.Attributes.UUID, downloadPath)
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("Sending new backup '%s' to configured destinations!", backup.Attributes.Name))
	errs := destinations.Backup(backupConfig.Destinations, backup.Attributes.UUID, backupFile)
	if len(errs) > 0 {
		log.Error(errs)
		return nil, fmt.Errorf("destination publishing failed with errors: %s", errs)
	}

	// Perform cleanup
	log.Debug(fmt.Sprintf("Deleting downloaded backup from temp path: %s", downloadPath))
	utils.DeleteFileIfExists(downloadPath)

	if deleteAfterBackup {
		log.Info(fmt.Sprintf("Deleting new backup '%s' from Pterodactyl server...", backup.Attributes.UUID))
		pterodactyl.DeleteServerBackup(*pterodactylServer, appServer, backup.Attributes.UUID)
	}

	return backup, nil
}

func startScheduler(configFile *config.Configuration, tmpDirPath string) {
	cronSchedular = cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))

	for _, backupConfig := range configFile.Backups {
		pterodactylServer, err := config.GetPterodactylServer(backupConfig.PterodactylServer, configFile.PterodactylServers)
		if err != nil {
			log.Fatal(err)
		}

		for _, appServerConfig := range backupConfig.Servers {
			appServer, err := pterodactyl.GetServer(*pterodactylServer, appServerConfig.Uuid)
			if err != nil {
				log.Fatal(err)
			}

			log.Info(fmt.Sprintf("Scheduling job for '%s' with schedule '%s'", appServer.Attributes.Name, appServerConfig.CronSchedule))

			cronSchedular.AddFunc(appServerConfig.CronSchedule, func() {
				_, err := PerformBackup(&backupConfig, pterodactylServer, appServer, false, tmpDirPath, appServerConfig.DeleteAfterBackup)
				utils.HandleError(err, false)
			})
		}
	}

	cronSchedular.Start()
}

func serve(configFilePath, tmpDirPath string) {
	configFile, err := config.ParseConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	startScheduler(configFile, tmpDirPath)

	log.Info("Server started")

	for {
		time.Sleep(time.Second)
	}
}

func Start(configFilePath, tmpDirPath string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Stopping application...")

		// Run Cleanup
		log.Info("Stopping schedular...")
		cronSchedular.Stop()

		os.Exit(0)
	}()

	serve(configFilePath, tmpDirPath)
}
