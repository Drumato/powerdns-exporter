package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	// TODO: メトリクス追加
	return reg
}
