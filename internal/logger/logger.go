package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func New() *slog.Logger {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug, 
		TimeFormat: time.Kitchen,    
		AddSource:  true,            
	})

	log := slog.New(handler)
	return log
}