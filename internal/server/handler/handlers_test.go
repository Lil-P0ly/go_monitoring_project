package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorageHandlerAddValue(t *testing.T) {
	ms := NewMSHandler()

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
			request := httptest.NewRequest(tt.method, "/", nil)

			request.SetPathValue("metrics_type", tt.path.metricType)
			request.SetPathValue("metrics_name", tt.path.metricName)
			request.SetPathValue("metrics_value", tt.path.metricValueStr)

			w := httptest.NewRecorder()

			ms.AddValue(w, request)
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
	ms := NewMSHandler()

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

			ms.NotFound(w, request)
			res := w.Result()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.wantResponse), strings.TrimSpace(string(resBody)))

		})
	}
}
