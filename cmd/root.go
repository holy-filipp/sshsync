package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		fmt.Println(viper.GetString("url"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func initConfig() error {
	viper.SetConfigName(CONFIG_FILENAME_NO_EXT)

	configPath, err := xdg.ConfigFile(CONFIG_DIR)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConfigInit, err)
	}

	slog.Info("config path defined", "path", configPath)

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		var fileNotFound viper.ConfigFileNotFoundError
		if errors.As(err, &fileNotFound) {
			slog.Warn("config file not found")
		} else {
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
