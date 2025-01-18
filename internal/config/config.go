package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bher20/pterodactyl-backup-manager/internal/destinations"
	"github.com/bher20/pterodactyl-backup-manager/internal/utils"
	"github.com/joho/godotenv"

	"github.com/bher20/pterodactyl-sdk-go/pkg/pterodactyl"
	log "github.com/sirupsen/logrus"
)

type Backup struct {
	PterodactylServer string                    `json:"pterodactylServer"`
	Servers           []Server                  `json:"servers"`
	Destinations      destinations.Destinations `json:"destinations"`
}

type Server struct {
	Uuid         string `json:"uuid"`
	Name         string `json:"name"`
	CronSchedule string `json:"cronSchedule"`
}

type Configuration struct {
	PterodactylServers []pterodactyl.PterodactylServer `json:"pterodactylServers"`
	Backups            []Backup                        `json:"backups"`
}

func setupEnv() {
	envFile := ".env"
	envFileFromEnv := os.Getenv("ENV_FILE")

	if envFileFromEnv != "" {
		envFile = envFileFromEnv
	}

	_, err := os.Stat(envFile)
	if err != nil {
		utils.HandleError(err, true)
	} else {
		err := godotenv.Load(envFile)
		utils.HandleError(err, true)
	}
}

func ParseConfig(configFile string) (*Configuration, error) {
	log.Info("Attempting to setup environment")
	setupEnv()

	log.Info(fmt.Sprintf("Attempting to read config file: %s", configFile))
	data, err := os.ReadFile(configFile)
	utils.HandleError(err, true)

	var configuration Configuration
	err = json.Unmarshal(data, &configuration)
	utils.HandleError(err, true)

	return &configuration, nil
}

func GetPterodatylServer(pterodactylServerName string, pterodactylServers []pterodactyl.PterodactylServer) (*pterodactyl.PterodactylServer, error) {
	var foundPterodactylServer *pterodactyl.PterodactylServer
	for _, pterodactylServer := range pterodactylServers {
		if pterodactylServer.Name == pterodactylServerName {
			foundPterodactylServer = &pterodactylServer
			break // Exit the loop once a match is found
		}
	}

	// Check if a matching struct was found
	if foundPterodactylServer != nil {
		return foundPterodactylServer, nil
	} else {
		return nil, fmt.Errorf("unable to find Pterodactyl server with name: %s", pterodactylServerName)
	}
}
