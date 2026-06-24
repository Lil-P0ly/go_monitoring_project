package models

import "strconv"

type MemStorage struct {
	Counter_metrics map[string][]int64
	Gauge_metrics   map[string][]float64
}

type Mestrics_type string

const (
	Metrics_type_counter Mestrics_type = "counter"
	Metrics_type_gauge   Mestrics_type = "gauge"
)

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Counter_metrics: make(map[string][]int64),
		Gauge_metrics:   make(map[string][]float64),
	}
}

func (ms *MemStorage) AddGauge(metric_name, metrics_value_str string) error {

	metrics_value, err := strconv.ParseFloat(metrics_value_str, 64)

	if err != nil {
		return err
	}
	ms.Gauge_metrics[metric_name] = append(ms.Gauge_metrics[metric_name], float64(metrics_value))
	return nil
}

func (ms *MemStorage) AddCounter(metric_name, metrics_value_str string) error {

	metrics_value, err := strconv.Atoi(metrics_value_str)

	if err != nil {
		return err
	}

	ms.Counter_metrics[metric_name] = append(ms.Counter_metrics[metric_name], int64(metrics_value))
	return nil
}
