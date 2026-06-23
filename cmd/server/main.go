package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type MemStorage struct {
	Counter_metrics map[string][]int64
	Gauge_metrics   map[string][]float64
}

func (ms *MemStorage) AddGauge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	metric_name := r.PathValue("metric_name")

	metrics_value_str := r.PathValue("metrics_value")

	metrics_value, err := strconv.ParseFloat(metrics_value_str, 64)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	ms.Gauge_metrics[metric_name] = append(ms.Gauge_metrics[metric_name], float64(metrics_value))

	w.WriteHeader(http.StatusOK)

}

func (ms *MemStorage) AddCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	metric_name := r.PathValue("metric_name")

	metrics_value_str := r.PathValue("metrics_value")

	metrics_value, err := strconv.Atoi(metrics_value_str)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	ms.Counter_metrics[metric_name] = append(ms.Counter_metrics[metric_name], int64(metrics_value))

	w.WriteHeader(http.StatusOK)

}

func (ms *MemStorage) PrintMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	response := ""
	response += fmt.Sprintf("Counter Metrics ---- \n")
	for k, v := range ms.Counter_metrics {
		response += fmt.Sprintf("Metric - %s: ", k)
		for _, val := range v {
			response += strconv.Itoa(int(val)) + " "
		}
		response += fmt.Sprintf("\n")
	}

	response += fmt.Sprintf("Gauge Metrics ---- \n")
	for k, v := range ms.Gauge_metrics {
		response += fmt.Sprintf("Metric - %s: ", k)
		for _, val := range v {
			response += fmt.Sprintf("%.2f", val) + " "
		}
		response += fmt.Sprintf("\n")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Counter_metrics: make(map[string][]int64),
		Gauge_metrics:   make(map[string][]float64),
	}
}

func main() {
	mux := http.NewServeMux()

	ms := NewMemStorage()
	mux.HandleFunc("/update/counter/{metric_name}/{metrics_value}", ms.AddCounter)
	mux.HandleFunc("/update/gauge/{metric_name}/{metrics_value}", ms.AddGauge)
	mux.HandleFunc("/metrics", ms.PrintMetrics)

	http.ListenAndServe(":8080", mux)

}
