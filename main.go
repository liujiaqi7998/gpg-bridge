package main

import (
	"fmt"
	"os"

	gpgbridgecmd "github.com/liujiaqi7998/gpg-bridge/cmd/gpg-bridge"
)

func main() {
	if err := gpgbridgecmd.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
