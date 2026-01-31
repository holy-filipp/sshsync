package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "this command fetches authorized keys from remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := url.ParseRequestURI(viper.GetString("url"))
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidUrl, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
