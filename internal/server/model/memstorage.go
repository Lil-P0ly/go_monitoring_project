package models

import "strconv"

type MemStorage struct {
	CounterMetrics map[string][]int64
	GaugeMetrics   map[string][]float64
}

type MetricsType string

const (
	MetricsTypeCounter MetricsType = "counter"
	MetricsTypeGauge   MetricsType = "gauge"
)

func NewMemStorage() *MemStorage {
	return &MemStorage{
		CounterMetrics: make(map[string][]int64),
		GaugeMetrics:   make(map[string][]float64),
	}
}

func (ms *MemStorage) AddGauge(metricName, metricValueStr string) error {

	metrics_value, err := strconv.ParseFloat(metricValueStr, 64)

	if err != nil {
		return err
	}
	ms.GaugeMetrics[metricName] = append(ms.GaugeMetrics[metricName], float64(metrics_value))
	return nil
}

func (ms *MemStorage) AddCounter(metricName, metricValueStr string) error {

	metrics_value, err := strconv.Atoi(metricValueStr)

	if err != nil {
		return err
	}

	ms.CounterMetrics[metricName] = append(ms.CounterMetrics[metricName], int64(metrics_value))
	return nil
}
