// package main for demoapp1 implements a single endpoint that fetches
// responses from demoapp2 and demoapp3 and combines them with its own message.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	appName     = "demoapp1"
	addrApp     = ":8080"
	addrMetrics = ":9000"

	reqTotalCollector = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: appName,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	statusSuccess     = string(http.StatusText(http.StatusOK))
	statusInternalErr = string(http.StatusText(http.StatusInternalServerError))
)

func main() {
	metricsRegistry := prometheus.NewRegistry()
	if err := metricsRegistry.Register(reqTotalCollector); err != nil {
		log.Fatal("Error registering Prometheus requests total collector:", err)
	}

	promHandler := promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{})

	appMux := http.NewServeMux()
	appMux.HandleFunc("/"+appName, handler)

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promHandler)

	appServer := &http.Server{
		Addr:    addrApp,
		Handler: appMux,
	}

	metricsServer := &http.Server{
		Addr:    addrMetrics,
		Handler: metricsMux,
	}

	go func() {
		log.Println("Starting demoapp1 server on:", addrApp)
		if err := appServer.ListenAndServe(); err != nil {
			log.Fatal("Error starting demoapp1 server:", err)
		}
	}()

	log.Println("Starting metrics server on:", addrMetrics)
	if err := metricsServer.ListenAndServe(); err != nil {
		log.Fatal("Error starting metrics server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", appName)

	demoApp2Addr := os.Getenv("DEMOAPP2_ADDR")
	demoApp3Addr := os.Getenv("DEMOAPP3_ADDR")

	demoApp2Payload, err := callDemoApp(demoApp2Addr + "/demoapp2")
	if err != nil {
		http.Error(w, "Error calling demoapp2: "+err.Error(), http.StatusInternalServerError)
		reqTotalCollector.WithLabelValues(r.URL.Path, r.Method, statusInternalErr).Inc()
		return
	}

	demoApp3Payload, err := callDemoApp(demoApp3Addr + "/demoapp3")
	if err != nil {
		http.Error(w, "Error calling demoapp3: "+err.Error(), http.StatusInternalServerError)
		reqTotalCollector.WithLabelValues(r.URL.Path, r.Method, statusInternalErr).Inc()
		return
	}

	fmt.Fprintf(w, "%s: %s %s", appName, string(demoApp2Payload), string(demoApp3Payload))

	reqTotalCollector.WithLabelValues(r.URL.Path, r.Method, statusSuccess).Inc()
}

func callDemoApp(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not make request to URL '%s': %w", url, err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read payload data for URL '%s': %w", url, err)
	}
	return data, nil
}
