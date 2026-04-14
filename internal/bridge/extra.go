package bridge

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/liujiaqi7998/gpg-bridge/internal/gpg"
	"github.com/liujiaqi7998/gpg-bridge/internal/protocol"
	"github.com/liujiaqi7998/gpg-bridge/internal/winpipe"
)

func BridgeExtra(ctx context.Context, listenAddr string, socketPath string) error {
	_ = gpg.PingAgent(ctx)

	listener, err := winpipe.Listen(listenAddr)
	if err != nil {
		return fmt.Errorf("listen extra bridge on %q: %w", listenAddr, err)
	}
	defer listener.Close()

	var (
		metaMu sync.Mutex
		cached *protocol.SocketMeta
	)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept extra connection: %w", err)
		}
		log.Printf("received extra request from remote: listen_addr=%q remote_addr=%q", listenAddr, conn.RemoteAddr())

		meta, err := func() (protocol.SocketMeta, error) {
			metaMu.Lock()
			defer metaMu.Unlock()
			if cached != nil {
				return *cached, nil
			}
			path := socketPath
			if path == "" {
				path, err = gpg.LoadSocketPath(ctx, gpg.SocketExtra)
				if err != nil {
					return protocol.SocketMeta{}, err
				}
			}
			loaded, err := gpg.LoadSocketMeta(path)
			if err != nil {
				return protocol.SocketMeta{}, err
			}
			cached = &loaded
			return loaded, nil
		}()
		if err != nil {
			_ = conn.Close()
			return fmt.Errorf("load extra socket metadata: %w", err)
		}

		go func(client net.Conn, socketMeta protocol.SocketMeta) {
			if err := delegateExtra(client, socketMeta); err != nil {
				metaMu.Lock()
				cached = nil
				metaMu.Unlock()
			}
		}(conn, meta)
	}
}

func delegateExtra(client net.Conn, meta protocol.SocketMeta) error {
	defer client.Close()

	target, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", meta.Port))
	if err != nil {
		_ = gpg.PingAgent(context.Background())
		return fmt.Errorf("connect target agent tcp: %w", err)
	}
	defer target.Close()

	if _, err := target.Write(meta.Nonce[:]); err != nil {
		return fmt.Errorf("write nonce: %w", err)
	}

	errCh := make(chan error, 2)
	go func() {
		_, err := io.Copy(target, client)
		if tcp, ok := target.(*net.TCPConn); ok {
			_ = tcp.CloseWrite()
		}
		errCh <- err
	}()
	go func() {
		_, err := io.Copy(client, target)
		if tcp, ok := client.(*net.TCPConn); ok {
			_ = tcp.CloseWrite()
		}
		errCh <- err
	}()

	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func loadSocketMetaFromFile(path string) (protocol.SocketMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return protocol.SocketMeta{}, err
	}
	return protocol.ParseSocketMeta(data)
}
