package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ErrInvalidUrl = errors.New("invalid url passed")
var ErrSetUrl = errors.New("failed to set url")

func setUrlInConfig(url string) error {
	viper.Set("url", url)

	configPath, err := xdg.ConfigFile(filepath.Join(CONFIG_DIR, CONFIG_FILENAME))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSetUrl, err)
	}

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("%w: %w", ErrSetUrl, err)
	}

	slog.Info("url set success")

	return nil
}

var urlCmd = &cobra.Command{
	Use:   "url [url to use]",
	Short: "set and save url to config file",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		_, err := url.ParseRequestURI(args[0])
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidUrl, err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return setUrlInConfig(args[0])
	},
}

func init() {
	rootCmd.AddCommand(urlCmd)
}
