package model

type Metrics interface {
	SendMetrics(url string) error
	CollectMetrics()
}
