package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"status", "route"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route"},
	)
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/healthz", handleHealthz)
	http.HandleFunc("/ready", handleReady)
	http.Handle("/metrics", promhttp.Handler())

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		httpRequestDuration.WithLabelValues("/").Observe(time.Since(start).Seconds())
	}()

	status := http.StatusOK
	defer func() {
		httpRequestsTotal.WithLabelValues(strconv.Itoa(status), "/").Inc()
	}()

	if shouldInjectFailure() && rand.Float32() < 0.1 {
		status = http.StatusInternalServerError
		http.Error(w, "Internal Server Error", status)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Progressive Delivery Demo</title></head>
<body>
<h1>Progressive Delivery Demo App</h1>
<p>Version: %s</p>
<p>Time: %s</p>
</body>
</html>`, getEnv("VERSION", "1.0.0"), time.Now().Format(time.RFC3339))
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		httpRequestDuration.WithLabelValues("/healthz").Observe(time.Since(start).Seconds())
	}()

	httpRequestsTotal.WithLabelValues("200", "/healthz").Inc()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		httpRequestDuration.WithLabelValues("/ready").Observe(time.Since(start).Seconds())
	}()

	httpRequestsTotal.WithLabelValues("200", "/ready").Inc()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func shouldInjectFailure() bool {
	return getEnv("INJECT_FAILURE", "false") == "true"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}