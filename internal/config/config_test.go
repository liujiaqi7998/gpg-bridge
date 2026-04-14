package config

import (
	"bytes"
	"testing"
)

func TestParseArgsRequiresAtLeastOneListener(t *testing.T) {
	var stderr bytes.Buffer

	_, err := ParseArgs(nil, &stderr)
	if err == nil {
		t.Fatal("expected error when no listeners are configured")
	}
}

func TestParseArgsParsesSupportedFlags(t *testing.T) {
	var stderr bytes.Buffer

	cfg, err := ParseArgs([]string{"--ssh", `\\.\\pipe\\gpg-bridge-ssh`, "--extra", "127.0.0.1:4321", "--extra-socket", `C:\\Users\\me\\AppData\\Roaming\\gnupg\\S.gpg-agent.extra`, "--detach"}, &stderr)
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
	if !cfg.Detach {
		t.Fatal("expected detach to be true")
	}
}
