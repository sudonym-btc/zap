package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/sudonym-btc/zap/cmd"
)

func main() {
	defer slog.MustClose()
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".zap")

	h := handler.MustFileHandler(configDir+"/debug.log", handler.WithLogLevels(slog.AllLevels))
	slog.PushHandler(h)
	slog.Std().Output = io.Discard

	cmd.Execute()

}
