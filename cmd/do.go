package cmd

import (
	"errors"
	"fmt"
	"github.com/holy-filipp/sshsync/lib"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ErrCantFetch = errors.New("can't fetch data from url")
var ErrRead = errors.New("failed to read body")
var ErrCantGetUser = errors.New("can't get user")
var ErrSshDirDoesntExist = errors.New("ssh dir doesn't exist")
var ErrCantWriteFile = errors.New("can't write file")

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "this command fetches authorized keys from remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCantGetUser, err)
		}

		sshDirPath := filepath.Join(usr.HomeDir, ".ssh")

		if _, err := os.Stat(sshDirPath); os.IsNotExist(err) {
			fmt.Printf("%s dir doesn't exist\n", sshDirPath)
			return ErrSshDirDoesntExist
		}

		u := viper.GetString("url")
		cacheBustUrl, err := url.Parse(u)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidUrl, err)
		}

		q := cacheBustUrl.Query()
		q.Set("kacache", lib.RandStringBytes(10))
		cacheBustUrl.RawQuery = q.Encode()

		resp, err := http.Get(cacheBustUrl.String())
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCantFetch, err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%w: response code is not 200: %d", ErrCantFetch, resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrRead, err)
		}

		if err := os.WriteFile(filepath.Join(sshDirPath, "authorized_keys"), data, 0600); err != nil {
			return fmt.Errorf("%w: %w", ErrCantWriteFile, err)
		}

		fmt.Println("authorized_keys synced")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
