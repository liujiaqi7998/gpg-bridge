package app

import (
	"context"
	"fmt"
	"log"

	"github.com/getlantern/systray"
	"github.com/liujiaqi7998/gpg-bridge/internal/bridge"
	"github.com/liujiaqi7998/gpg-bridge/internal/config"
)

type App struct {
	cfg config.Config
}

func New(cfg config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := newRunResult()
	go func() {
		result.finish(bridge.Run(ctx, a.cfg))
	}()

	onReady := func() {
		if icon, err := loadTrayIcon(); err != nil {
			log.Printf("load tray icon failed: %v", err)
		} else {
			systray.SetIcon(icon)
		}
		systray.SetTitle("gpg-bridge")
		systray.SetTooltip("gpg-bridge is running")
		quitItem := systray.AddMenuItem("Quit", "Quit gpg-bridge")
		go func() {
			<-quitItem.ClickedCh
			log.Print("tray quit requested")
			cancel()
			result.finish(nil)
			systray.Quit()
		}()
	}

	onExit := func() {}
	go func() {
		err := result.wait()
		if err != nil {
			log.Printf("bridge stopped with error: %v", err)
		} else {
			log.Print("bridge stopped")
		}
		systray.Quit()
	}()

	systray.Run(onReady, onExit)
	if err := result.wait(); err != nil && ctx.Err() == nil {
		return fmt.Errorf("run bridge: %w", err)
	}
	return nil
}
