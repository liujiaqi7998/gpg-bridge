package main

import (
	"fmt"
	"log"
	"os"

	gpgbridgecmd "github.com/liujiaqi7998/gpg-bridge/cmd/gpg-bridge"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stderr)
	log.Println("gpg-bridge starting...")
	if err := gpgbridgecmd.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
