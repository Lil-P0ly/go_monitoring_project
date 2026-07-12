package models

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGauge(t *testing.T) {

	var tests = []struct {
		name     string
		input    MemStorage
		want     MemStorage
		addvalue float64
		addkey   string
	}{
		{
			name:     "Test-1",
			input:    MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			want:     MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1, 1.2}}},
			addvalue: 1.2,
			addkey:   "bar",
		},

		{
			name:     "Test-2",
			input:    MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			want:     MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}, "fuz": {-0.999}}},
			addvalue: -0.999,
			addkey:   "fuz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans_ms := tt.input
			ans_ms.AddGauge(tt.addkey, tt.addvalue)
			assert.Equal(t, ans_ms, tt.want)

		})
	}
}

func TestAddCounter(t *testing.T) {

	var tests = []struct {
		name     string
		input    MemStorage
		want     MemStorage
		addvalue int64
		addkey   string
	}{
		{
			name:     "Test-1",
			input:    MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			want:     MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2, 12}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			addvalue: 12,
			addkey:   "bar",
		},

		{
			name:     "Test-2",
			input:    MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			want:     MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2}, "fuz": {-2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			addvalue: -2,
			addkey:   "fuz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans_ms := tt.input
			ans_ms.AddCounter(tt.addkey, tt.addvalue)
			assert.Equal(t, ans_ms, tt.want)

		})
	}
}

func TestGetGauges(t *testing.T) {

	var tests = []struct {
		name  string
		input MemStorage
		want  map[string][]float64
	}{
		{
			name:  "Test-1",
			input: MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {2}}, GaugeMetrics: map[string][]float64{"foo": {1.0, 1.01}, "bar": {1.1}}},
			want:  map[string][]float64{"foo": {1.0, 1.01}, "bar": {1.1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans_map := tt.input.GetGauges()
			assert.Equal(t, ans_map, tt.want)

		})
	}
}

func TestGetCounters(t *testing.T) {

	var tests = []struct {
		name  string
		input MemStorage
		want  map[string][]int64
	}{
		{
			name:  "Test-1",
			input: MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {1, 2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			want:  map[string][]int64{"foo": {1}, "bar": {1, 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans_map := tt.input.GetCounters()
			assert.Equal(t, ans_map, tt.want)

		})
	}
}

func TestGetLastGauge(t *testing.T) {

	var tests = []struct {
		name             string
		ms               MemStorage
		metricName       string
		wantValue        float64
		wantError        bool
		wantErrorMessage string
	}{
		{
			name:       "Test-Exist-Metric",
			ms:         MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {1, 2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			metricName: "bar",
			wantValue:  1.1,
			wantError:  false,
		},
		{
			name:             "Test-NonExist-Metric",
			ms:               MemStorage{CounterMetrics: map[string][]int64{"foo": {1}, "bar": {1, 2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			metricName:       "foof",
			wantError:        true,
			wantErrorMessage: "Metric not found in MemStorage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.ms.GetLastGauge(tt.metricName)

			if !tt.wantError {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, val)
			} else {
				log.Println(err)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrorMessage)

			}

		})
	}
}

func TestGetLastCounter(t *testing.T) {

	var tests = []struct {
		name             string
		ms               MemStorage
		metricName       string
		wantValue        int64
		wantError        bool
		wantErrorMessage string
	}{
		{
			name:       "Test-Exist-Metric",
			ms:         MemStorage{CounterMetrics: map[string][]int64{"du": {1}, "mb": {1, 2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			metricName: "mb",
			wantValue:  2,
			wantError:  false,
		},
		{
			name:             "Test-NonExist-Metric",
			ms:               MemStorage{CounterMetrics: map[string][]int64{"du": {1}, "mb": {1, 2}}, GaugeMetrics: map[string][]float64{"foo": {1.0}, "bar": {1.1}}},
			metricName:       "foof",
			wantError:        true,
			wantErrorMessage: "Metric not found in MemStorage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.ms.GetLastCounter(tt.metricName)

			if !tt.wantError {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, val)
			} else {
				log.Println(err)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrorMessage)

			}

		})
	}
}
