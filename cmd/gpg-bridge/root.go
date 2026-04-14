package gpgbridgecmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/liujiaqi7998/gpg-bridge/internal/bridge"
	"github.com/liujiaqi7998/gpg-bridge/internal/config"
)

func Run(args []string) error {
	cfg, err := config.ParseArgs(args, os.Stderr)
	if err != nil {
		return err
	}
	if cfg.Detach {
		return detach(args)
	}
	return bridge.Run(context.Background(), cfg)
}

func detach(args []string) error {
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "--detach" {
			filtered = append(filtered, arg)
		}
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable: %w", err)
	}
	cmd := exec.Command(exe, filtered...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000200 | 0x00000008 | 0x04000000}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start detached process: %w", err)
	}
	return nil
}

func Usage() string {
	return strings.TrimSpace(`gpg-bridge bridges OpenSSH and GnuPG agent endpoints on Windows.`)
}
