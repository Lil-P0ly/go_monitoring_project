package models

import "errors"

type MemStorage struct {
	CounterMetrics map[string]int64
	GaugeMetrics   map[string][]float64
}

type MetricsType string

const (
	MetricsTypeCounter MetricsType = "counter"
	MetricsTypeGauge   MetricsType = "gauge"
)

func NewMemStorage() *MemStorage {
	return &MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string][]float64),
	}
}

func (ms *MemStorage) AddGauge(metricName string, metricValue float64) {

	ms.GaugeMetrics[metricName] = append(ms.GaugeMetrics[metricName], metricValue)
}

func (ms *MemStorage) AddCounter(metricName string, metricValue int64) {

	_, ok := ms.CounterMetrics[metricName]
	if !ok {
		ms.CounterMetrics[metricName] = metricValue
	} else {
		ms.CounterMetrics[metricName] += metricValue
	}
}

func (ms *MemStorage) GetGauges() map[string][]float64 {
	return ms.GaugeMetrics
}

func (ms *MemStorage) GetCounters() map[string]int64 {
	return ms.CounterMetrics
}

func (ms *MemStorage) GetLastGauge(metricName string) (float64, error) {
	slice, ok := ms.GaugeMetrics[metricName]
	if !ok {
		return 0, errors.New("Metric not found in MemStorage")
	}
	if len(slice) == 0 {
		return 0, errors.New("Metrics-slice is empty")
	}
	return slice[len(slice)-1], nil
}

func (ms *MemStorage) GetLastCounter(metricName string) (int64, error) {
	val, ok := ms.CounterMetrics[metricName]
	if !ok {
		return 0, errors.New("Metric not found in MemStorage")
	}

	return val, nil
}
