package server

import (
	"log/slog"
	"net/http"

	"github.com/Drumato/powerdns-exporter/middleware"
	powerdns "github.com/Drumato/powerdns-exporter/powerdns/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(
	logger *slog.Logger,
	client powerdns.Client,
	reg *prometheus.Registry,
	address string,
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := client.Healthcheck(r.Context()); err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", middleware.MetricsProvider(logger, client, promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})))

	return &http.Server{
		Handler: mux,
		Addr:    address,
	}
}
