package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	gpgbridgecmd "github.com/liujiaqi7998/gpg-bridge/cmd/gpg-bridge"
	"github.com/liujiaqi7998/gpg-bridge/internal/logx"
	"github.com/liujiaqi7998/gpg-bridge/internal/singleinstance"
)

func main() {
	hideConsoleWindow()

	guard, err := singleinstance.Acquire()
	if err != nil {
		if errors.Is(err, singleinstance.ErrAlreadyActive) {
			showInfoMessage("gpg-bridge", "Already Acttive")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer guard.Close()

	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logCloser, err := logx.Configure(filepath.Clean(baseDir))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer logCloser.Close()

	if err := gpgbridgecmd.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
