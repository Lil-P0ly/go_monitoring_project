package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
	"github.com/go-chi/chi/v5"
)

type MemoryStorageHandler struct {
	Storage models.Storage
}

func NewMSHandlerWithStorage(s models.Storage) *MemoryStorageHandler {
	return &MemoryStorageHandler{Storage: s}
}

func (msh *MemoryStorageHandler) AddValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "metrics_type")
	metricName := chi.URLParam(r, "metrics_name")
	metricValueStr := chi.URLParam(r, "metrics_value")

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

func (msh *MemoryStorageHandler) PrintLastValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "metrics_type")
	metricName := chi.URLParam(r, "metrics_name")

	switch metricType {
	case string(models.MetricsTypeGauge):
		log.Println("get last gauge value")
		val, err := msh.Storage.GetLastGauge(metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		str := fmt.Sprintf("%f", val)
		w.Write([]byte(str))

	case string(models.MetricsTypeCounter):
		log.Println("get last counter value")
		val, err := msh.Storage.GetLastCounter(metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		str := fmt.Sprintf("%d", val)
		w.Write([]byte(str))

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

func (msh *MemoryStorageHandler) PrintMetricsHTML(w http.ResponseWriter, r *http.Request) {

	type ViewData struct {
		Title          string
		GaugeMetrics   map[string][]float64
		CounterMetrics map[string][]int64
	}
	data := ViewData{
		Title:          "Metrics Project",
		GaugeMetrics:   msh.Storage.GetGauges(),
		CounterMetrics: msh.Storage.GetCounters(),
	}
	tmpl, err := template.ParseFiles("/home/user/dev/go/yp/go-pro/go_monitoring_project/internal/server/templates/index.html")

	if err != nil {
		log.Println("Fail to open template fail for main page")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
