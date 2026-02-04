package common

import (
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/gorilla/mux"
)

const metricsName = "petstore_http_requests_total"

type requestKey struct {
	Path   string
	Method string
}

type MetricsCounter struct {
	mutex         sync.Mutex
	requestCounts map[requestKey]uint64
}

func PromMetricsCounter() *MetricsCounter {
	return &MetricsCounter{
		requestCounts: make(map[requestKey]uint64),
	}
}

func (m *MetricsCounter) Increment(path, method string) {
	m.mutex.Lock()
	m.requestCounts[requestKey{Path: path, Method: method}]++
	m.mutex.Unlock()
}

func (m *MetricsCounter) snapshot() map[requestKey]uint64 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	snap := make(map[requestKey]uint64, len(m.requestCounts))
	for k, v := range m.requestCounts {
		snap[k] = v
	}
	return snap
}

func (m *MetricsCounter) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snap := m.snapshot()

		// Sort keys for deterministic output.
		keys := make([]requestKey, 0, len(snap))
		for k := range snap {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			if keys[i].Path != keys[j].Path {
				return keys[i].Path < keys[j].Path
			}
			return keys[i].Method < keys[j].Method
		})

		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		fmt.Fprintf(w, "# HELP %s Total number of HTTP requests.\n", metricsName)
		fmt.Fprintf(w, "# TYPE %s counter\n", metricsName)
		for _, k := range keys {
			fmt.Fprintf(w, "%s{path=%q,method=%q} %d\n", metricsName, k.Path, k.Method, snap[k])
		}
	})
}

func (m *MetricsCounter) Middleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			path := "unknown"
			route := mux.CurrentRoute(r)
			if route != nil {
				tmpl, err := route.GetPathTemplate()
				if err == nil {
					path = tmpl
				}
			}
			m.Increment(path, r.Method)
		})
	}
}
