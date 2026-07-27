package model

type MetricSender interface {
	SendMetrics(string) error
	SendMetricJSON(string) error
	CollectMetrics()
}