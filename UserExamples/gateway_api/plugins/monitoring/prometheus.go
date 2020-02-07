package monitoring

import (
	p "github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// MetricsWrapper prometheus
func MetricsWrapper(h http.Handler) http.Handler {
	ph := p.Handler()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			ph.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}
