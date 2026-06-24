package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
)

type MemoryStorageHandler struct {
	Storage *models.MemStorage
}

func NewMSHandler() *MemoryStorageHandler {
	return &MemoryStorageHandler{
		Storage: models.NewMemStorage(),
	}
}

func (msh *MemoryStorageHandler) AddValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metricType := r.PathValue("metrics_type")
	metricName := r.PathValue("metrics_name")
	metricValueStr := r.PathValue("metrics_value")

	switch metricType {
	case string(models.MetricsTypeGauge):
		log.Println("update gauge")
		metricValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		msh.Storage.AddGauge(metricName, metricValue)

	case string(models.MetricsTypeCounter):
		log.Println("update counter")

		metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		msh.Storage.AddCounter(metricName, metricValue)

	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (msh *MemoryStorageHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	log.Println("Not found handler")
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metricType := r.PathValue("metrics_type")

	if metricType != string(models.MetricsTypeCounter) && metricType != string(models.MetricsTypeGauge) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)

}

func (msh *MemoryStorageHandler) PrintMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	response := ""
	response += fmt.Sprintf("Counter Metrics ---- \n")
	for k, v := range msh.Storage.CounterMetrics {
		response += fmt.Sprintf("Metric - %s: ", k)
		for _, val := range v {
			response += strconv.Itoa(int(val)) + " "
		}
		response += fmt.Sprintf("\n")
	}
	response += fmt.Sprintf("Gauge Metrics ---- \n")
	for k, v := range msh.Storage.GaugeMetrics {
		response += fmt.Sprintf("Metric - %s: ", k)
		for _, val := range v {
			response += fmt.Sprintf("%.2f ", val)
		}
		response += fmt.Sprintf("\n")
	}
	w.Write([]byte(response))
}
