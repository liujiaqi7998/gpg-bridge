package logx

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenLogFileRotatesLargeFile(t *testing.T) {
	tmpDir := t.TempDir()
	logDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	logPath := filepath.Join(logDir, "gpg-bridge.log")
	large := make([]byte, maxLogSizeBytes+128)
	if err := os.WriteFile(logPath, large, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	file, err := openLogFile(logPath)
	if err != nil {
		t.Fatalf("openLogFile returned error: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		t.Fatalf("Stat returned error: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("expected rotated log file size 0, got %d", info.Size())
	}
}

func TestConfigureCreatesLogDirectoryAndWritesLog(t *testing.T) {
	tmpDir := t.TempDir()
	oldWriter := log.Writer()
	oldFlags := log.Flags()
	oldPrefix := log.Prefix()
	defer log.SetOutput(oldWriter)
	defer log.SetFlags(oldFlags)
	defer log.SetPrefix(oldPrefix)

	closer, err := Configure(tmpDir)
	if err != nil {
		t.Fatalf("Configure returned error: %v", err)
	}
	defer closer.Close()

	log.Print("hello log file")

	data, err := os.ReadFile(filepath.Join(tmpDir, "logs", "gpg-bridge.log"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected log file to contain data")
	}
}
