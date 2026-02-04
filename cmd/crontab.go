package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var ErrExecutingCommand = errors.New("failed to execute command")
var ErrGetExecutable = errors.New("failed to get executable")
var ErrCrontab = errors.New("failed to write crontab")

var crontabCmd = &cobra.Command{
	Use:   "crontab [optional url]",
	Short: "add sync job to user's cron",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}

		_, err := url.ParseRequestURI(args[0])
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidUrl, err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			if err := setUrlInConfig(args[0]); err != nil {
				return fmt.Errorf("%w: %w", ErrSetUrl, err)
			}
		}

		progPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrGetExecutable, err)
		}

		progDir := filepath.Dir(progPath)

		command := fmt.Sprintf("* * * * * %s do >> %s 2>&1\n", progPath, filepath.Join(progDir, "sshsync_crontab.log"))

		currentCron, err := exec.Command("/bin/sh", "-c", "crontab -l").CombinedOutput()
		if err != nil && !strings.Contains(string(currentCron), "no crontab for") {
			return fmt.Errorf("%w: %s: %w: %s", ErrExecutingCommand, "crontab -l", err, string(currentCron))
		}

		if strings.Contains(string(currentCron), command) {
			slog.Info("crontab job already exists")
			return nil
		}

		var cron strings.Builder
		if !strings.Contains(string(currentCron), "no crontab for") {
			cron.Write(currentCron)
		}

		cron.WriteString(command)

		writeCronCmd := exec.Command("/bin/sh", "-c", "crontab -")
		writeCronCmd.Stdin = strings.NewReader(cron.String())

		writeCronOut, err := writeCronCmd.Output()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrGetExecutable, err)
		}

		if string(writeCronOut) != "" {
			return fmt.Errorf("%w: %s", ErrCrontab, writeCronOut)
		}

		slog.Info("added job to user's crontab")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(crontabCmd)
}
