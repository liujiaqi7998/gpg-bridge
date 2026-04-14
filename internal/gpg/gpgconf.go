package gpg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/liujiaqi7998/gpg-bridge/internal/protocol"
)

type SocketName string

const (
	SocketSSH   SocketName = "agent-ssh-socket"
	SocketExtra SocketName = "agent-extra-socket"
)

func LoadSocketPath(ctx context.Context, name SocketName) (string, error) {
	cmd := exec.CommandContext(ctx, "gpgconf", "--list-dir", string(name))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("run gpgconf for %s: %w: %s", name, err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func PingAgent(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "gpg-connect-agent", "/bye")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("run gpg-connect-agent: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func LoadSocketMeta(path string) (protocol.SocketMeta, error) {
	if _, err := os.Stat(path); err != nil {
		return protocol.SocketMeta{}, err
	}
	data, err := os.ReadFile(strings.ReplaceAll(path, `\`, `/`))
	if err != nil {
		return protocol.SocketMeta{}, err
	}
	return protocol.ParseSocketMeta(data)
}
