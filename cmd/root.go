/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bherville/pterodactyl-backup-manager/internal/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	tmpDirPath    string
	debugLogging  bool
	traceLoggging bool
)

type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}
func setLogging(cmd *cobra.Command, args []string) {
	if debugLogging {
		log.Info("Debug logs enabled")
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
	} else if traceLoggging {
		log.Info("Trace logs enabled")
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&log.TextFormatter{
			DisableLevelTruncation: true,
		})
	} else {
		plainFormatter := new(PlainFormatter)
		log.SetFormatter(plainFormatter)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pterodactyl-backup-manager",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setLogging(cmd, args)

		err := utils.CreateDirIfNotExists(tmpDirPath, 0755)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

// PersistentPreRunE: func(cmd *Command, args []string): {
// 	err := utils.CreateDirIfNotExists(tmpDirPath, 0755)
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}

// 	return err
// }

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Logging flags
	rootCmd.PersistentFlags().BoolVarP(&debugLogging, "debug", "d", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&traceLoggging, "trace", "t", false, "enable trace logging")

	// Config flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.json", "config file")

	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	rootCmd.PersistentFlags().StringVar(&tmpDirPath, "tmp-dir", filepath.Join(dir, "tmp"), "config file")
}
