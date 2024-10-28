package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	powerdns "github.com/Drumato/powerdns-exporter/powerdns/v1"
)

func MetricsProvider(logger *slog.Logger, client powerdns.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.DebugContext(r.Context(), "start updating metrics")
		servers, err := client.GetServers(r.Context())
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to get PowerDNS Servers", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, s := range servers {
			fmt.Println(s)
		}

		next.ServeHTTP(w, r)
	})
}
