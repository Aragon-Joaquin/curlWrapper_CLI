package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	fileName = "logs.log"
)

var logDir string

// executes automatically when starting the program
func init() {
	wd, err := os.Getwd()

	if err != nil {
		slog.Warn("failed to get user cache dir; falling back to temp dir", "err", err, "path", wd)
	}

	logDir = wd
}

func DefaultPath() string {
	return filepath.Join(logDir, fileName)
}

func Load(path string, level slog.Level) error {
	_, err := os.Stat(filepath.Join(logDir, fileName))

	if os.IsNotExist(err) {
		os.Create(fileName)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewTextHandler(file, opts)
	slog.SetDefault(slog.New(handler))
	return nil
}
