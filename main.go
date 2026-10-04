package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var operationDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "calculator_operation_duration_seconds",
		Help: "Time spent executing calculator operations.",
	},
	[]string{"operation"},
)

func init() {
	prometheus.MustRegister(operationDuration)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		fmt.Println("Metrics server running on http://localhost:8080/metrics")
		http.ListenAndServe(":8080", nil)
	}()

	for {
		var a, b float64
		var operation string

		fmt.Print("Enter first number: ")
		fmt.Scan(&a)

		fmt.Print("Enter operation (+, -, *, /): ")
		fmt.Scan(&operation)

		fmt.Print("Enter second number: ")
		fmt.Scan(&b)

		start := time.Now()

		var result float64

		switch operation {
		case "+":
			result = a + b
		case "-":
			result = a - b
		case "*":
			result = a * b
		case "/":
			result = a / b
		default:
			fmt.Println("Invalid operation")
			continue
		}

		duration := time.Since(start)

		operationDuration.WithLabelValues(operation).Observe(duration.Seconds())

		fmt.Println("Result:", result)
		fmt.Println("Operation duration:", duration)
	}
}
