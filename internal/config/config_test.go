package config

import (
	"bytes"
	"testing"
)

func TestParseArgsUsesDefaultExtraWhenNoArgs(t *testing.T) {
	var stderr bytes.Buffer

	cfg, err := ParseArgs(nil, &stderr)
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}
	if cfg.Extra != "127.0.0.1:35132" {
		t.Fatalf("unexpected default extra address: %q", cfg.Extra)
	}
	if cfg.SSH != "" {
		t.Fatalf("unexpected default ssh address: %q", cfg.SSH)
	}
}

func TestParseArgsParsesSupportedFlags(t *testing.T) {
	var stderr bytes.Buffer

	cfg, err := ParseArgs([]string{"--ssh", `\\.\\pipe\\gpg-bridge-ssh`, "--extra", "127.0.0.1:4321", "--extra-socket", `C:\\Users\\me\\AppData\\Roaming\\gnupg\\S.gpg-agent.extra`}, &stderr)
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if cfg.SSH != `\\.\\pipe\\gpg-bridge-ssh` {
		t.Fatalf("unexpected ssh address: %q", cfg.SSH)
	}
	if cfg.Extra != "127.0.0.1:4321" {
		t.Fatalf("unexpected extra address: %q", cfg.Extra)
	}
	if cfg.ExtraSocket != `C:\\Users\\me\\AppData\\Roaming\\gnupg\\S.gpg-agent.extra` {
		t.Fatalf("unexpected extra socket: %q", cfg.ExtraSocket)
	}
}

func TestParseArgsRejectsUnexpectedPositionalArguments(t *testing.T) {
	var stderr bytes.Buffer

	_, err := ParseArgs([]string{"unexpected"}, &stderr)
	if err == nil {
		t.Fatal("expected positional arguments to fail")
	}
}
