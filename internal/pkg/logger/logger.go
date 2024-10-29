package logger

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"path/filepath"
	"time"
)

func SetupLogger(config *config.Config) *zerolog.Logger {
	logFileName := fmt.Sprintf("request_log_%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(config.LogFilePath, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Error opening log file")
	}

	// Create the logger
	zerolog.New(logFile).With().Timestamp().Caller().Logger()
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.New(logFile).With().Timestamp().Caller().Logger()

	log := zerolog.New(logFile).With().Caller().Timestamp().Logger()
	zerolog.DefaultContextLogger = &log

	return &log
}

func GetFromContext(ctx context.Context) *zerolog.Logger {
	if logger, ok := ctx.Value(types.LoggerContextKey).(*zerolog.Logger); ok {
		return logger
	} else {
		return zerolog.DefaultContextLogger
	}
}
