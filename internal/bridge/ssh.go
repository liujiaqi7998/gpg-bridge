package bridge

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/liujiaqi7998/gpg-bridge/internal/pageant"
	"github.com/liujiaqi7998/gpg-bridge/internal/winpipe"
)

func BridgeSSH(ctx context.Context, listenAddr string) error {
	_ = ctx
	listener, err := winpipe.Listen(listenAddr)
	if err != nil {
		return fmt.Errorf("listen ssh bridge on %q: %w", listenAddr, err)
	}
	defer listener.Close()

	sem := make(chan struct{}, 4)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept ssh connection: %w", err)
		}
		sem <- struct{}{}
		go func(c net.Conn) {
			defer func() { <-sem }()
			_ = handleSSHConn(c)
		}(conn)
	}
}

func handleSSHConn(conn net.Conn) error {
	defer conn.Close()

	for {
		var lenBuf [4]byte
		if _, err := io.ReadFull(conn, lenBuf[:]); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil
			}
			return err
		}
		msgLen := binary.BigEndian.Uint32(lenBuf[:])
		payload := make([]byte, int(msgLen)+4)
		copy(payload[:4], lenBuf[:])
		if _, err := io.ReadFull(conn, payload[4:]); err != nil {
			return err
		}
		resp, err := pageant.Query(payload)
		if err != nil {
			return err
		}
		if _, err := conn.Write(resp); err != nil {
			return err
		}
	}
}
