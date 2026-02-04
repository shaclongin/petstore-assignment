package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestMetricsCounter_Increment(t *testing.T) {
	mc := PromMetricsCounter()

	mc.Increment("/api/pet", "GET")
	mc.Increment("/api/pet", "GET")
	mc.Increment("/api/pet", "POST")
	mc.Increment("/api/pet/{petId}", "GET")

	snap := mc.snapshot()
	require.Equal(t, uint64(2), snap[requestKey{Path: "/api/pet", Method: "GET"}])
	require.Equal(t, uint64(1), snap[requestKey{Path: "/api/pet", Method: "POST"}])
	require.Equal(t, uint64(1), snap[requestKey{Path: "/api/pet/{petId}", Method: "GET"}])
}

func TestMetricsCounter_Handler(t *testing.T) {
	mc := PromMetricsCounter()
	mc.Increment("/api/pet", "GET")
	mc.Increment("/api/pet", "POST")
	mc.Increment("/api/pet/{petId}", "GET")

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	mc.Handler().ServeHTTP(w, req)

	body := w.Body.String()
	require.Equal(t, "text/plain; version=0.0.4; charset=utf-8", w.Header().Get("Content-Type"))
	require.Contains(t, body, "# HELP petstore_http_requests_total Total number of HTTP requests.")
	require.Contains(t, body, "# TYPE petstore_http_requests_total counter")
	require.Contains(t, body, `petstore_http_requests_total{path="/api/pet",method="GET"} 1`)
	require.Contains(t, body, `petstore_http_requests_total{path="/api/pet",method="POST"} 1`)
	require.Contains(t, body, `petstore_http_requests_total{path="/api/pet/{petId}",method="GET"} 1`)
}

func TestMetricsCounter_Middleware(t *testing.T) {
	mc := PromMetricsCounter()

	r := mux.NewRouter()
	r.HandleFunc("/api/pet", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET", "POST")
	r.HandleFunc("/api/pet/{petId}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
	r.Use(mc.Middleware())

	// Make requests
	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/api/pet"},
		{"GET", "/api/pet"},
		{"POST", "/api/pet"},
		{"GET", "/api/pet/123"},
		{"GET", "/api/pet/456"},
	}
	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	snap := mc.snapshot()
	require.Equal(t, uint64(2), snap[requestKey{Path: "/api/pet", Method: "GET"}])
	require.Equal(t, uint64(1), snap[requestKey{Path: "/api/pet", Method: "POST"}])
	require.Equal(t, uint64(2), snap[requestKey{Path: "/api/pet/{petId}", Method: "GET"}])
}

func TestMetricsCounter_Concurrent(t *testing.T) {
	mc := PromMetricsCounter()
	done := make(chan struct{})

	// run 10 coroutines concurrently
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				mc.Increment("/api/pet", "GET")
			}
			done <- struct{}{}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	snap := mc.snapshot()
	require.Equal(t, uint64(1000), snap[requestKey{Path: "/api/pet", Method: "GET"}])
}
