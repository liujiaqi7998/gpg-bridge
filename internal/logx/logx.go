package logx

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const maxLogSizeBytes int64 = 1 << 20
const logDirName = "logs"
const logFileName = "gpg-bridge.log"

type closerFunc func() error

func (fn closerFunc) Close() error {
	return fn()
}

func Configure(baseDir string) (io.Closer, error) {
	logDir := filepath.Join(baseDir, logDirName)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	file, err := openLogFile(filepath.Join(logDir, logFileName))
	if err != nil {
		return nil, err
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(file)
	return closerFunc(file.Close), nil
}

func openLogFile(path string) (*os.File, error) {
	info, err := os.Stat(path)
	if err == nil && info.Size() >= maxLogSizeBytes {
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("remove oversized log file: %w", err)
		}
	}
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("stat log file: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	return file, nil
}
