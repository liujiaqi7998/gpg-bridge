package config

import (
	"flag"
	"fmt"
	"io"
)

const DefaultExtraAddr = "127.0.0.1:35132"

type Config struct {
	SSH         string
	Extra       string
	ExtraSocket string
	Detach      bool
}

func ParseArgs(args []string, stderr io.Writer) (Config, error) {
	fs := flag.NewFlagSet("gpg-bridge", flag.ContinueOnError)
	fs.SetOutput(stderr)

	cfg := Config{Extra: DefaultExtraAddr}
	fs.StringVar(&cfg.SSH, "ssh", "", "Sets the listenning address to bridge the ssh socket")
	fs.StringVar(&cfg.Extra, "extra", DefaultExtraAddr, "Sets the listenning to bridge the extra socket")
	fs.StringVar(&cfg.ExtraSocket, "extra-socket", "", "Sets the path to gnupg extra socket optionaly")
	fs.BoolVar(&cfg.Detach, "detach", false, "Runs the program as a background daemon")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage of %s:\n", fs.Name())
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	if fs.NArg() > 0 {
		return Config{}, fmt.Errorf("unexpected positional arguments: %v", fs.Args())
	}
	return cfg, nil
}
