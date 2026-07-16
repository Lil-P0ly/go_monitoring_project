package model

type Metrics interface {
	SendMetrics(string) error
	CollectMetrics()
}