package handler

import (
	_ "embed"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/Lil-P0ly/go_monitoring_project/internal/server/logger"
	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type MemoryStorageHandler struct {
	Storage models.Storage
}

//go:embed templates/index.html
var indexTmpl string

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
		logger.Info("update gauge")
		metricValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		msh.Storage.AddGauge(metricName, metricValue)

	case string(models.MetricsTypeCounter):
		logger.Info("update counter")
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
		logger.Info("get last gauge value")
		val, err := msh.Storage.GetLastGauge(metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))

	case string(models.MetricsTypeCounter):
		logger.Info("get last counter value")
		val, err := msh.Storage.GetLastCounter(metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatInt(val, 10)))

	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}

func (msh *MemoryStorageHandler) PrintMetricsHTML(w http.ResponseWriter, r *http.Request) {

	type ViewData struct {
		Title          string
		GaugeMetrics   map[string][]float64
		CounterMetrics map[string]int64
	}
	data := ViewData{
		Title:          "Metrics Project",
		GaugeMetrics:   msh.Storage.GetGauges(),
		CounterMetrics: msh.Storage.GetCounters(),
	}
	tmpl, err := template.New("index").Parse(indexTmpl)

	if err != nil {
		logger.Error("Fail to open template fail for main page")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func (msh *MemoryStorageHandler) LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		duration := time.Since(startTime)

		logger.Info("Request",
			zap.String("Method", r.Method),
			zap.String("URL", r.URL.Path),
			zap.Duration("Duration", duration),
		)

		logger.Info("Response",
			zap.Int("StatusCode", lrw.statusCode),
			zap.Int("Size", lrw.size),
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}
