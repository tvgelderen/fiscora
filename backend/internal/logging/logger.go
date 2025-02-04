package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/tvgelderen/fiscora/internal/config"
)

type RotatingFileHandler struct {
	mu           sync.Mutex
	currentFile  *os.File
	baseFilePath string
	maxSize      int64
	currentSize  int64
	handler      slog.Handler
}

func NewRotatingFileHandler(basePath string, maxSizeMB int) (*RotatingFileHandler, error) {
	if maxSizeMB <= 0 {
		return nil, fmt.Errorf("maxSizeMB must be positive")
	}

	file, err := openOrCreateFile(basePath)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	h := &RotatingFileHandler{
		currentFile:  file,
		baseFilePath: basePath,
		maxSize:      int64(maxSizeMB) * 1024 * 1024,
		currentSize:  info.Size(),
	}

	h.handler = slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{
					Key:   "timestamp",
					Value: a.Value,
				}
			}
			return a
		},
	})

	return h, nil
}

func (h *RotatingFileHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	estimatedSize := int64(len(r.Message) + 100)

	if h.currentSize+estimatedSize > h.maxSize {
		if err := h.rotate(); err != nil {
			return err
		}
	}

	err := h.handler.Handle(ctx, r)
	if err == nil {
		h.currentSize += estimatedSize
	}

	return err
}

func (h *RotatingFileHandler) rotate() error {
	if err := h.currentFile.Close(); err != nil {
		return err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	dir, file := filepath.Split(h.baseFilePath)
	ext := filepath.Ext(file)
	base := file[:len(file)-len(ext)]
	rotatedPath := filepath.Join(dir, fmt.Sprintf("%s-%s%s", base, timestamp, ext))

	if err := os.Rename(h.baseFilePath, rotatedPath); err != nil {
		return err
	}

	newFile, err := openOrCreateFile(h.baseFilePath)
	if err != nil {
		return err
	}

	h.currentFile = newFile
	h.currentSize = 0
	h.handler = slog.NewJSONHandler(newFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{
					Key:   "timestamp",
					Value: a.Value,
				}
			}
			return a
		},
	})

	return nil
}

func openOrCreateFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("log directory does not exist: %s. Please run the setup script first", dir)
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

func (h *RotatingFileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &RotatingFileHandler{
		currentFile:  h.currentFile,
		baseFilePath: h.baseFilePath,
		maxSize:      h.maxSize,
		currentSize:  h.currentSize,
		handler:      h.handler.WithAttrs(attrs),
	}
}

func (h *RotatingFileHandler) WithGroup(name string) slog.Handler {
	return &RotatingFileHandler{
		currentFile:  h.currentFile,
		baseFilePath: h.baseFilePath,
		maxSize:      h.maxSize,
		currentSize:  h.currentSize,
		handler:      h.handler.WithGroup(name),
	}
}

func (h *RotatingFileHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func SetupLogger() (*slog.Logger, error) {
	if !config.Env.Production {
		return slog.New(slog.NewTextHandler(os.Stdout, nil)), nil
	}

	logDir := "/var/log/fiscora-backend"
	logFile := filepath.Join(logDir, "app.log")

	handler, err := NewRotatingFileHandler(logFile, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to create rotating handler: %v", err)
	}
	return slog.New(handler), nil
}
