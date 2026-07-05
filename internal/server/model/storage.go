package models

type Storage interface {
	AddGauge(metricName string, metricValue float64)
	AddCounter(metricName string, metricValue int64)

	GetGauges() map[string][]float64
	GetCounters() map[string][]int64
}
