package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(reg *prometheus.Registry, address string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// TODO: PowerDNSとつながることをテスト
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	return &http.Server{
		Handler: mux,
		Addr:    address,
	}
}
