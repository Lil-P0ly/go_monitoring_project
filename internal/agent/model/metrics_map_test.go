package model

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMetricsToStruct(t *testing.T) {

	m := NewMetricsMap()

	// type metricsSample struct {
	// 	Alloc  uint64
	// 	NextGC uint64
	// }

	tests := []struct {
		name           string
		metricsSample  runtime.MemStats
		wantMetricsMap MetricsMap
	}{
		{
			name:           "Test-1",
			metricsSample:  runtime.MemStats{Alloc: 1, NextGC: 3},
			wantMetricsMap: MetricsMap{GaugeMetrics: map[string]float64{"Alloc": 1.0, "NextGC": 3.0}, CounterMetrics: map[string]int64{"PollCount": 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.ParseMetricsToStruct(tt.metricsSample)
			assert.Equal(t, m.CounterMetrics, tt.wantMetricsMap.CounterMetrics)
			assert.Equal(t, m.GaugeMetrics["Alloc"], tt.wantMetricsMap.GaugeMetrics["Alloc"])

		})
	}
}
