package bridge

import (
	"context"
	"errors"
	"sync"

	"github.com/liujiaqi7998/gpg-bridge/internal/config"
)

func Run(ctx context.Context, cfg config.Config) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	started := 0

	if cfg.Extra != "" {
		started++
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- BridgeExtra(ctx, cfg.Extra, cfg.ExtraSocket)
		}()
	}
	if cfg.SSH != "" {
		started++
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- BridgeSSH(ctx, cfg.SSH)
		}()
	}
	if started == 0 {
		return errors.New("no bridge configured")
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}
