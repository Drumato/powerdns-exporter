package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/Drumato/powerdns-exporter/metrics"
	powerdns "github.com/Drumato/powerdns-exporter/powerdns/v1"
	"github.com/Drumato/powerdns-exporter/server"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	Logger *slog.Logger
	Client powerdns.Client

	host                    string
	port                    string
	powerDNSWebserverHost   string
	powerDNSWebserverPort   string
	powerDNSWebserverScheme string
}

func New(logger *slog.Logger, apiKey string) *cobra.Command {
	command := RootCommand{
		Logger: logger,
	}

	c := &cobra.Command{
		Use: "powerdns-exporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			powerDNSWebserverAddr := net.JoinHostPort(command.powerDNSWebserverHost, command.powerDNSWebserverPort)
			baseURL := fmt.Sprintf("%s://%s", command.powerDNSWebserverScheme, powerDNSWebserverAddr)
			command.Client = powerdns.NewDefault(command.Logger, baseURL, apiKey)
			if err := command.Run(cmd.Context()); err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
	}

	c.Flags().StringVar(&command.host, "host", "", "exporter host")
	c.Flags().StringVar(&command.port, "port", "8080", "exporter port")
	c.Flags().StringVar(&command.powerDNSWebserverHost, "powerdns-webserver-host", "pdns", "the PowerDNS webserver host")
	c.Flags().StringVar(&command.powerDNSWebserverPort, "powerdns-webserver-port", "8081", "the PowerDNS webserver port")
	c.Flags().StringVar(&command.powerDNSWebserverScheme, "powerdns-webserver-scheme", "http", "the PowerDNS webserver protocol scheme")
	return c
}

func (c *RootCommand) Run(ctx context.Context) error {
	addr := net.JoinHostPort(c.host, c.port)
	reg := metrics.NewRegistry()
	server := server.New(c.Logger, c.Client, reg, addr)

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
