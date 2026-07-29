package handler

import (
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func (msh *MemoryStorageHandler) AddValueJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	var metricsModel models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metricsModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metricsModel.MType {
	case string(models.MetricsTypeGauge):
		logger.Info("update gauge")
		if metricsModel.Value == nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		msh.Storage.AddGauge(metricsModel.ID, *metricsModel.Value)
		val, err := msh.Storage.GetLastGauge(metricsModel.ID)
		if err != nil {
			http.Error(w, "Internal error while update json-gauge response", http.StatusInternalServerError)
			return
		}
		metricsModel.Value = &val

	case string(models.MetricsTypeCounter):
		logger.Info("update counter")
		if metricsModel.Delta == nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		msh.Storage.AddCounter(metricsModel.ID, *metricsModel.Delta)
		delta, err := msh.Storage.GetLastCounter(metricsModel.ID)
		if err != nil {
			http.Error(w, "Internal error while update json-counter response", http.StatusInternalServerError)
			return
		}
		metricsModel.Delta = &delta

	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(metricsModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (msh *MemoryStorageHandler) GetValueJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	var metricsModel models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricsModel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metricsModel.MType {
	case string(models.MetricsTypeGauge):
		val, err := msh.Storage.GetLastGauge(metricsModel.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		metricsModel.Value = &val

	case string(models.MetricsTypeCounter):
		delta, err := msh.Storage.GetLastCounter(metricsModel.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		metricsModel.Delta = &delta

	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metricsModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func (msh *MemoryStorageHandler) GzipDecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "gzip" {
			next.ServeHTTP(w, r)
			return
		}

		if ct := r.Header.Get("Content-Type"); ct != "application/json" && ct != "text/html" {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "invalid gzip body", http.StatusBadRequest)
			return
		}
		defer gz.Close()

		r.Body = io.NopCloser(gz)

		next.ServeHTTP(w, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (msh *MemoryStorageHandler) GzipСompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		if ct := r.Header.Get("Content-Type"); ct != "application/json" && ct != "text/html" {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
