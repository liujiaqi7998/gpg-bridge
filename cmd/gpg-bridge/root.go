package gpgbridgecmd

import (
	"fmt"
	"os"

	"github.com/liujiaqi7998/gpg-bridge/internal/app"
	"github.com/liujiaqi7998/gpg-bridge/internal/config"
)

func Run(args []string) error {
	cfg, err := config.ParseArgs(args, os.Stderr)
	if err != nil {
		return err
	}
	return app.New(cfg).Run()
}

func Usage() string {
	return "gpg-bridge runs in the system tray and bridges OpenSSH and GnuPG on Windows."
}

func formatAlreadyRunningError() error {
	return fmt.Errorf("gpg-bridge is already running")
}
