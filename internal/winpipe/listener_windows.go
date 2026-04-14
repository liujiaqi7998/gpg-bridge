//go:build windows

package winpipe

import (
	"fmt"
	"net"
	"strings"

	"github.com/Microsoft/go-winio"
)

func Listen(addr string) (net.Listener, error) {
	if strings.HasPrefix(addr, `\\.\pipe\`) {
		ln, err := winio.ListenPipe(addr, nil)
		if err != nil {
			return nil, fmt.Errorf("listen named pipe: %w", err)
		}
		return ln, nil
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen tcp: %w", err)
	}
	return ln, nil
}
