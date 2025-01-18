/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bher20/pterodactyl-backup-manager/internal/config"
	"github.com/bher20/pterodactyl-backup-manager/internal/server"
	"github.com/bher20/pterodactyl-backup-manager/internal/utils"
	"github.com/bher20/pterodactyl-sdk-go/pkg/pterodactyl"

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
	PreRun: setLogging,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := config.ParseConfig(cfgFile)
		utils.HandleError(err, true)

		for _, backupConfig := range configFile.Backups {
			pterodatylServer, err := config.GetPterodatylServer(backupConfig.PterodactylServer, configFile.PterodactylServers)
			utils.HandleError(err, true)

			for _, appServerConfig := range backupConfig.Servers {
				appServer, err := pterodactyl.GetServer(*pterodatylServer, appServerConfig.Uuid)
				utils.HandleError(err, true)

				utils.HandleError(server.PerformBackup(&backupConfig, pterodatylServer, appServer, true), true)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
