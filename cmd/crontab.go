package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrExecutingCommand = errors.New("failed to execute command")
var ErrGetExecutable = errors.New("failed to get executable")
var ErrCrontab = errors.New("failed to write crontab")

var crontabCmd = &cobra.Command{
	Use:   "crontab",
	Short: "add sync job to user's cron",
	Args:  validateUrlArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		progPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrGetExecutable, err)
		}

		progDir := filepath.Dir(progPath)

		currentCron, err := exec.Command("/bin/sh", "-c", "crontab -l").Output()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrExecutingCommand, err)
		}

		var cron strings.Builder
		if !strings.HasPrefix(string(currentCron), "no crontab for") {
			cron.Write(currentCron)
		}

		cron.WriteString(fmt.Sprintf("\n* * * * * %s do >> %s 2>&1", progPath, filepath.Join(progDir, "sshsync_crontab.log")))

		writeCronCmd := exec.Command("/bin/sh", "-c", "crontab -")
		writeCronCmd.Stdin = strings.NewReader(cron.String())

		writeCronOut, err := writeCronCmd.Output()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrGetExecutable, err)
		}

		if string(writeCronOut) != "" {
			return fmt.Errorf("%w: %s", ErrCrontab, writeCronOut)
		}

		if err := setUrlInConfig(args[0]); err != nil {
			return fmt.Errorf("%w: %w", ErrSetUrl, err)
		}

		fmt.Println("added sshsync to user's crontab")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(crontabCmd)
}
