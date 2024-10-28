package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/Drumato/powerdns-exporter/cmd"
)

const (
	powerDNSAPIKeyEnv = "POWERDNS_API_KEY"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var logLevel slog.Level

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))

	apiKey, ok := os.LookupEnv(powerDNSAPIKeyEnv)
	if !ok {
		logger.ErrorContext(ctx, "the environment variable 'POWERDNS_API_KEY' must be defined")
		os.Exit(1)
	}

	if err := cmd.New(logger, apiKey).ExecuteContext(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to run", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
