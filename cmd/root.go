package cmd

import (
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

const CONFIG_DIR = "sshsync"
const CONFIG_FILENAME = "config.json"
const CONFIG_FILENAME_NO_EXT = "config"

var ErrInit = errors.New("init failed")
var ErrConfigInit = errors.New("failed to initialize config")

var rootCmd = &cobra.Command{
	Use:   "sshsync",
	Short: "sshsync is a tool for syncing ssh authorized_keys on multiple machines",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info(fmt.Sprintf("sshsync uses URL %s", viper.GetString("url")))
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() error {
	viper.SetConfigName(CONFIG_FILENAME_NO_EXT)

	configPath, err := xdg.ConfigFile(CONFIG_DIR)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConfigInit, err)
	}

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		var fileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &fileNotFound) {
			return fmt.Errorf("%w: %w", ErrConfigInit, err)
		}
	}

	return nil
}

func init() {
	if err := initConfig(); err != nil {
		slog.Error(fmt.Errorf("%w: %w", ErrInit, err).Error())
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringP("url", "u", "", "url to fetch authorized_keys from")

	if err := viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url")); err != nil {
		slog.Error(fmt.Errorf("%w: %w", ErrInit, err).Error())
		os.Exit(1)
	}
}
