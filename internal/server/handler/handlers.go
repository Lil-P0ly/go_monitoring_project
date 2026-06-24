package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
)

type MememoryStorageHandler struct {
	Storage *models.MemStorage
}

func NewMSHandler() *MememoryStorageHandler {
	return &MememoryStorageHandler{
		Storage: models.NewMemStorage(),
	}
}

func (msh *MememoryStorageHandler) AddValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metric_type := r.PathValue("metrics_type")
	metric_name := r.PathValue("metric_name")
	metrics_value_str := r.PathValue("metrics_value")

	switch metric_type {
	case string(models.Metrics_type_gauge):
		log.Println("update gauge")
		err := msh.Storage.AddGauge(metric_name, metrics_value_str)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	case string(models.Metrics_type_counter):
		log.Println("update counter")

		err := msh.Storage.AddCounter(metric_name, metrics_value_str)
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

func (msh *MememoryStorageHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	log.Println("Not found handler")
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metric_type := r.PathValue("metrics_type")

	if metric_type != string(models.Metrics_type_counter) && metric_type != string(models.Metrics_type_gauge) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	http.Error(w, "Not Found", http.StatusNotFound)

}

func (msh *MememoryStorageHandler) PrintMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	response := ""
	response += fmt.Sprintf("Counter Metrics ---- \n")
	for k, v := range msh.Storage.Counter_metrics {
		response += fmt.Sprintf("Metric - %s: ", k)
		for _, val := range v {
			response += strconv.Itoa(int(val)) + " "
		}
		response += fmt.Sprintf("\n")
	}
	response += fmt.Sprintf("Gauge Metrics ---- \n")
	for k, v := range msh.Storage.Gauge_metrics {
		response += fmt.Sprintf("Metric - %s: ", k)
		for _, val := range v {
			response += fmt.Sprintf("%.2f ", val)
		}
		response += fmt.Sprintf("\n")
	}
	w.Write([]byte(response))
}
