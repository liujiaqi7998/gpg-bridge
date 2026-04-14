//go:build windows

package winpipe

import (
	"fmt"
	"net"
	"strings"
	"time"

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
	tcpLn, ok := ln.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("listen tcp returned unexpected listener type %T", ln)
	}
	return &tcpListener{TCPListener: tcpLn}, nil
}

type tcpListener struct {
	*net.TCPListener
}

func (l *tcpListener) Accept() (net.Conn, error) {
	for {
		if err := l.SetDeadline(time.Now().Add(500 * time.Millisecond)); err != nil {
			return nil, err
		}
		conn, err := l.AcceptTCP()
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			continue
		}
		return conn, err
	}
}
