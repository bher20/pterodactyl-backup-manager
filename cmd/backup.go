/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bherville/pterodactyl-backup-manager/internal/config"
	"github.com/bherville/pterodactyl-backup-manager/internal/server"
	"github.com/bherville/pterodactyl-backup-manager/internal/utils"
	"github.com/bherville/pterodactyl-sdk-go/pkg/pterodactyl"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := config.ParseConfig(cfgFile)
		utils.HandleError(err, true)

		for _, backupConfig := range configFile.Backups {
			pterodactylServer, err := config.GetPterodactylServer(backupConfig.PterodactylServer, configFile.PterodactylServers)
			utils.HandleError(err, true)

			for _, appServerConfig := range backupConfig.Servers {
				appServer, err := pterodactyl.GetServer(*pterodactylServer, appServerConfig.Uuid)
				utils.HandleError(err, true)

				_, err = server.PerformBackup(&backupConfig, pterodactylServer, appServer, true, tmpDirPath, appServerConfig.DeleteAfterBackup)
				utils.HandleError(err, true)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
