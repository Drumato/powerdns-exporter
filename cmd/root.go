package cmd

import (
	"context"
	"log/slog"
	"net"

	"github.com/Drumato/powerdns-exporter/metrics"
	"github.com/Drumato/powerdns-exporter/server"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	Logger *slog.Logger

	host string
	port string
}

func New(logger *slog.Logger) *cobra.Command {
	command := RootCommand{
		Logger: logger,
	}

	c := &cobra.Command{
		Use: "powerdns-exporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := command.Run(cmd.Context()); err != nil {
				return errors.WithStack(err)
			}
			return nil
		},
	}

	c.Flags().StringVar(&command.host, "host", "", "exporter host")
	c.Flags().StringVar(&command.port, "port", "8080", "exporter port")
	return c
}

func (c *RootCommand) Run(ctx context.Context) error {
	addr := net.JoinHostPort(c.host, c.port)
	reg := metrics.NewRegistry()
	server := server.New(reg, addr)

	c.Logger.InfoContext(ctx, "starting exporter...")
	go func(ctx context.Context) {
		if err := server.ListenAndServe(); err != nil {
			c.Logger.ErrorContext(ctx, "failed to listen and serve", slog.String("error", err.Error()))
		}
	}(ctx)

	<-ctx.Done()
	c.Logger.InfoContext(ctx, "signal interupped; shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		c.Logger.ErrorContext(ctx, "failed to shutdown server", slog.String("error", err.Error()))
	}

	return nil
}
