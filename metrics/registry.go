package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func NewRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.Register(collectors.NewGoCollector())
	// TODO: メトリクス追加
	return reg
}
