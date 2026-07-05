package models

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

func (ms *MemStorage) AddGauge(metricName string, metricValue float64) {

	ms.GaugeMetrics[metricName] = append(ms.GaugeMetrics[metricName], metricValue)
}

func (ms *MemStorage) AddCounter(metricName string, metricValue int64) {

	ms.CounterMetrics[metricName] = append(ms.CounterMetrics[metricName], metricValue)
}

func (ms *MemStorage) GetGauges() map[string][]float64 {
	return ms.GaugeMetrics
}

func (ms *MemStorage) GetCounters() map[string][]int64 {
	return ms.CounterMetrics
}
