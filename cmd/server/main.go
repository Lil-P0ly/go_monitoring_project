package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type MemStorage struct {
	Counter_metrics map[string][]int64
	Gauge_metrics   map[string][]float64
}

type Mestrics_type string

const (
	Metrics_type_counter Mestrics_type = "counter"
	Metrics_type_gauge   Mestrics_type = "gauge"
)

func (ms *MemStorage) AddGauge(metric_name, metrics_value_str string) error {

	metrics_value, err := strconv.ParseFloat(metrics_value_str, 64)

	if err != nil {
		return err
	}
	ms.Gauge_metrics[metric_name] = append(ms.Gauge_metrics[metric_name], float64(metrics_value))
	return nil
}

func (ms *MemStorage) AddCounter(metric_name, metrics_value_str string) error {

	metrics_value, err := strconv.Atoi(metrics_value_str)

	if err != nil {
		return err
	}

	ms.Counter_metrics[metric_name] = append(ms.Counter_metrics[metric_name], int64(metrics_value))
	return nil
}

func (ms *MemStorage) AddValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metric_type := r.PathValue("metrics_type")
	metric_name := r.PathValue("metric_name")
	metrics_value_str := r.PathValue("metrics_value")

	switch metric_type {
	case string(Metrics_type_gauge):
		log.Println("update gauge")
		err := ms.AddGauge(metric_name, metrics_value_str)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	case string(Metrics_type_counter):
		log.Println("update counter")

		err := ms.AddCounter(metric_name, metrics_value_str)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ms *MemStorage) NotFound(w http.ResponseWriter, r *http.Request) {
	log.Println("Not found handler")
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metric_type := r.PathValue("metrics_type")

	if metric_type != string(Metrics_type_counter) && metric_type != string(Metrics_type_gauge) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	http.Error(w, "Not Found", http.StatusNotFound)

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

	mux.HandleFunc("/update/{metrics_type}/{metric_name}/{metrics_value}", ms.AddValue)

	mux.HandleFunc("/update/{metrics_type}/", ms.NotFound)

	mux.HandleFunc("/metrics", ms.PrintMetrics)

	http.ListenAndServe(":8080", mux)

}
