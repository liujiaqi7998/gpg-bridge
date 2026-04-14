package config

import (
	"flag"
	"fmt"
	"io"
)

type Config struct {
	SSH         string
	Extra       string
	ExtraSocket string
	Detach      bool
}

func ParseArgs(args []string, stderr io.Writer) (Config, error) {
	fs := flag.NewFlagSet("gpg-bridge", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var cfg Config
	fs.StringVar(&cfg.SSH, "ssh", "", "Sets the listenning address to bridge the ssh socket")
	fs.StringVar(&cfg.Extra, "extra", "", "Sets the listenning to bridge the extra socket")
	fs.StringVar(&cfg.ExtraSocket, "extra-socket", "", "Sets the path to gnupg extra socket optionaly")
	fs.BoolVar(&cfg.Detach, "detach", false, "Runs the program as a background daemon")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage of %s:\n", fs.Name())
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	if cfg.SSH == "" && cfg.Extra == "" {
		return Config{}, fmt.Errorf("at least one of --ssh or --extra is required")
	}
	if fs.NArg() > 0 {
		return Config{}, fmt.Errorf("unexpected positional arguments: %v", fs.Args())
	}
	return cfg, nil
}
