package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "get version",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("0.0.3-pre-release")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
