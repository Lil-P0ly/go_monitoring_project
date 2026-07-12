package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorageHandlerAddValue(t *testing.T) {
	memStorage := models.NewMemStorage()

	msh := NewMSHandlerWithStorage(memStorage)

	r := chi.NewRouter()
	r.HandleFunc("/update/{metrics_type}/{metrics_name}/{metrics_value}", msh.AddValue)

	type path struct {
		metricType     string
		metricName     string
		metricValueStr string
	}
	tests := []struct {
		name           string
		method         string
		path           path
		wantStatusCode int
		wantResponse   string
	}{
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			path:           path{metricType: "gauge", metricName: "vault", metricValueStr: "11.1"},
			wantStatusCode: http.StatusMethodNotAllowed,
			wantResponse:   "Method Not Allowed",
		},

		{
			name:   "Invalid metrics type",
			method: http.MethodPost,
			path:   path{metricType: "bool", metricName: "vault", metricValueStr: "11.1"},

			wantStatusCode: http.StatusBadRequest,
			wantResponse:   "Bad Request",
		},
		{
			name:           "Invalid gauge value",
			method:         http.MethodPost,
			path:           path{metricType: "gauge", metricName: "vault", metricValueStr: "x"},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   "Bad Request",
		},
		{
			name:           "Invalid counter value",
			method:         http.MethodPost,
			path:           path{metricType: "counter", metricName: "cnt", metricValueStr: "x"},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   "Bad Request",
		},
		{
			name:           "Valid gauge value",
			method:         http.MethodPost,
			path:           path{metricType: "gauge", metricName: "vault", metricValueStr: "11.1"},
			wantStatusCode: http.StatusOK,
			wantResponse:   "",
		},
		{
			name:           "Valid counter value",
			method:         http.MethodPost,
			path:           path{metricType: "counter", metricName: "cnt", metricValueStr: "11"},
			wantStatusCode: http.StatusOK,
			wantResponse:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			path_str := fmt.Sprintf("/update/%s/%s/%s", tt.path.metricType, tt.path.metricName, tt.path.metricValueStr)
			request := httptest.NewRequest(tt.method, path_str, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.wantResponse), strings.TrimSpace(string(resBody)))

		})
	}
}

func TestNotFound(t *testing.T) {
	memStorage := models.NewMemStorage()

	msh := NewMSHandlerWithStorage(memStorage)
	type path struct {
		metricType string
		metricName string
	}
	tests := []struct {
		name           string
		method         string
		path           path
		wantStatusCode int
		wantResponse   string
	}{
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			path:           path{metricType: "gauge", metricName: "vault"},
			wantStatusCode: http.StatusMethodNotAllowed,
			wantResponse:   "Method Not Allowed",
		},

		{
			name:   "Invalid metrics type",
			method: http.MethodPost,
			path:   path{metricType: "bool", metricName: "vault"},

			wantStatusCode: http.StatusBadRequest,
			wantResponse:   "Bad Request",
		},
		{
			name:           "Not Found method",
			method:         http.MethodPost,
			path:           path{metricType: "gauge", metricName: "vault"},
			wantStatusCode: http.StatusNotFound,
			wantResponse:   "Not Found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", nil)

			request.SetPathValue("metrics_type", tt.path.metricType)
			request.SetPathValue("metrics_name", tt.path.metricName)

			w := httptest.NewRecorder()

			msh.NotFound(w, request)
			res := w.Result()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.wantResponse), strings.TrimSpace(string(resBody)))

		})
	}
}

func TestMemoryStorageHandlerPrintLastValue(t *testing.T) {
	memStorage := models.NewMemStorage()

	msh := NewMSHandlerWithStorage(memStorage)

	r := chi.NewRouter()
	r.Handle("/update/{metrics_type}/{metrics_name}", http.HandlerFunc(msh.PrintLastValue))

	type path struct {
		metricType string
		metricName string
	}
	tests := []struct {
		name           string
		method         string
		path           path
		prep           func()
		wantStatusCode int
		wantResponse   string
	}{
		{
			name:           "Invalid method",
			method:         http.MethodPost,
			path:           path{metricType: "gauge", metricName: "vault"},
			wantStatusCode: http.StatusMethodNotAllowed,
			wantResponse:   "Method Not Allowed",
		},
		{
			name:           "Invalid metrics type",
			method:         http.MethodGet,
			path:           path{metricType: "bool", metricName: "vault"},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   "Bad Request",
		},
		{
			name:   "Valid gauge value",
			method: http.MethodGet,
			path:   path{metricType: "gauge", metricName: "vault"},
			prep: func() {
				msh.Storage.AddGauge("vault", 11.1)
			},
			wantStatusCode: http.StatusOK,
			wantResponse:   "11.100000",
		},
		{
			name:   "Valid counter value",
			method: http.MethodGet,
			path:   path{metricType: "counter", metricName: "cnt"},
			prep: func() {
				msh.Storage.AddCounter("cnt", 42)
			},
			wantStatusCode: http.StatusOK,
			wantResponse:   "42",
		},
		{
			name:           "Not found gauge",
			method:         http.MethodGet,
			path:           path{metricType: "gauge", metricName: "missing"},
			wantStatusCode: http.StatusNotFound,
			wantResponse:   "Metric not fount in MemStorage",
		},
		{
			name:           "Not found counter",
			method:         http.MethodGet,
			path:           path{metricType: "counter", metricName: "missing"},
			wantStatusCode: http.StatusNotFound,
			wantResponse:   "Metric not fount in MemStorage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prep != nil {
				tt.prep()
			}

			pathStr := fmt.Sprintf("/update/%s/%s", tt.path.metricType, tt.path.metricName)
			request := httptest.NewRequest(tt.method, pathStr, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.wantResponse), strings.TrimSpace(string(resBody)))
		})
	}
}

func TestMemoryStorageHandlerPrintMetricsHTML(t *testing.T) {
	t.Run("Render with metrics", func(t *testing.T) {
		memStorage := models.NewMemStorage()
		msh := NewMSHandlerWithStorage(memStorage)

		msh.Storage.AddGauge("vault", 11.1)
		msh.Storage.AddCounter("cnt", 42)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		msh.PrintMetricsHTML(w, request)
		res := w.Result()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		bodyStr := string(body)
		assert.Contains(t, bodyStr, "Metrics Project")
		assert.Contains(t, bodyStr, "vault")
		assert.Contains(t, bodyStr, "cnt")
		assert.Contains(t, bodyStr, "11.1")
		assert.Contains(t, bodyStr, "42")
	})

	t.Run("Render empty storage", func(t *testing.T) {
		memStorage := models.NewMemStorage()
		msh := NewMSHandlerWithStorage(memStorage)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		msh.PrintMetricsHTML(w, request)
		res := w.Result()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Contains(t, string(body), "Metrics Project")
	})
}
