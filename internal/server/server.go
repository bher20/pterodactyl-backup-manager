package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bher20/pterodactyl-backup-manager/internal/config"
	"github.com/bher20/pterodactyl-backup-manager/internal/destinations"
	"github.com/bher20/pterodactyl-backup-manager/internal/utils"

	"github.com/bher20/pterodactyl-sdk-go/pkg/pterodactyl"

	"github.com/robfig/cron"

	log "github.com/sirupsen/logrus"
)

var (
	cronSchedular *cron.Cron
)

const (
	WAIT_FOR_BACKUP_SECONDS = 2
)

func PerformBackup(backupConfig *config.Backup, pterodatylServer *pterodactyl.PterodactylServer, appServer pterodactyl.Server, errorsFatal bool) error {
	log.Info(fmt.Sprintf("Running backup of '%s'...", appServer.Attributes.Name))

	backup, err := pterodactyl.BackupServer(*pterodatylServer, appServer)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Backup '%s' started!", backup.Attributes.UUID))

	// Wait until backup is completed on the pterodatylServer side
	log.Info(fmt.Sprintf("Waiting for '%s' to finish...", backup.Attributes.UUID))
	for {
		backup, err = pterodactyl.GetServerBackup(*pterodatylServer, appServer, backup.Attributes.UUID)
		log.Trace(fmt.Sprintf("Backup '%s' backup.Attributes.CompletedAt is zero: %t with err: %s!", backup.Attributes.UUID, time.Time.IsZero(backup.Attributes.CompletedAt), err))

		if !time.Time.IsZero(backup.Attributes.CompletedAt) {
			log.Info(fmt.Sprintf("Backup '%s' was completed at '%s'!", backup.Attributes.UUID, backup.Attributes.CompletedAt.Format(utils.TIME_FORMAT)))
			time.Sleep(2 * time.Second)
			break
		}
		log.Trace(fmt.Sprintf("Still waiting for '%s' to finish...", backup.Attributes.UUID))
	}

	log.Info(fmt.Sprintf("Downlading new backup '%s'!", backup.Attributes.UUID))
	backupFileBytes, err := pterodactyl.DownloadServerBackup(*pterodatylServer, appServer, backup.Attributes.UUID)
	if err != nil {
		return err
	}

	if len(backupFileBytes) == 0 {
		return fmt.Errorf("backup size is 0 bytes")
	}

	log.Info(fmt.Sprintf("Sending new backup '%s' to configured destinations!", backup.Attributes.Name))
	errs := destinations.Backup(backupConfig.Destinations, backup.Attributes.UUID, backupFileBytes)
	if len(errs) > 0 {
		log.Error(errs)
		return fmt.Errorf("destination publishing failed with errors: %s", errs)
	}

	return nil
}

func startScheduler(configFile *config.Configuration) {
	cronSchedular = cron.New()

	for _, backupConfig := range configFile.Backups {
		pterodatylServer, err := config.GetPterodatylServer(backupConfig.PterodactylServer, configFile.PterodactylServers)
		if err != nil {
			log.Fatal(err)
		}

		for _, appServerConfig := range backupConfig.Servers {
			appServer, err := pterodactyl.GetServer(*pterodatylServer, appServerConfig.Uuid)
			if err != nil {
				log.Fatal(err)
			}

			log.Info(fmt.Sprintf("Scheduling job for %s with schedule '%s'", appServer.Attributes.Name, appServerConfig.CronSchedule))

			cronSchedular.AddFunc(appServerConfig.CronSchedule, func() {
				PerformBackup(&backupConfig, pterodatylServer, appServer, false)
			})
		}
	}

	cronSchedular.Start()
}

func serve(configFilePath string) {
	configFile, err := config.ParseConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	startScheduler(configFile)

	for {
		time.Sleep(time.Second)
	}
}

func Start(configFilePath string) {
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

	serve(configFilePath)
}
