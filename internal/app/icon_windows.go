//go:build windows

package app

import (
	"fmt"
	"os"
	"path/filepath"
)

func loadTrayIcon() ([]byte, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("resolve executable path: %w", err)
	}
	iconPath := filepath.Join(filepath.Dir(exePath), "icon.ico")
	data, err := os.ReadFile(iconPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", iconPath, err)
	}
	return data, nil
}
