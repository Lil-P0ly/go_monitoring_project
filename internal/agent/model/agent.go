package model

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

var MetricsNames = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second

type MetricsMap struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

func (m *MetricsMap) SendMetrics(url string) error {

	log.Println("Start send metric to server")
	for metric_name, metric_value := range m.GaugeMetrics {
		url := fmt.Sprintf("http://%s/update/gauge/%v/%f", url, metric_name, metric_value)
		_, err := http.Post(url, "text/plain", nil)

		if err != nil {
			log.Printf("Failed send metrics to server, %v", err)
			return err
		}
		// log.Println(resp.StatusCode)
	}

	for metric_name, metric_value := range m.CounterMetrics {
		url := fmt.Sprintf("http://%s/update/counter/%v/%d", url, metric_name, metric_value)
		_, err := http.Post(url, "text/plain", nil)

		if err != nil {
			log.Printf("Failed send metrics to server, %v", err)
			return err
		}
		// log.Println(resp.StatusCode)
	}
	return nil

}

func (m *MetricsMap) CollectMetrics() {

	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	s := reflect.ValueOf(memStat)
	var val float64

	log.Println("Load Runtime metrics")
	for _, metric := range MetricsNames {
		field := s.FieldByName(metric)

		switch field.Kind() {
		case reflect.Uint64, reflect.Uint32, reflect.Uint, reflect.Uint16, reflect.Uint8:
			val = float64(field.Uint())

		case reflect.Int64, reflect.Int32, reflect.Int, reflect.Int16, reflect.Int8:
			val = float64(field.Int())

		case reflect.Float32, reflect.Float64:
			val = float64(field.Float())

		default:
			log.Printf("unsupported field type: %s", field.Kind())
			continue
		}
		m.GaugeMetrics[metric] = val
	}

	if mapVal, ok := m.CounterMetrics["PollCount"]; ok {
		mapVal++
		m.CounterMetrics["PollCount"] = mapVal
	} else {
		m.CounterMetrics["PollCount"] = 1
	}
	m.GaugeMetrics["RandomValue"] = randFloat(-100000, 100000)

}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func NewMetricsMap() *MetricsMap {
	return &MetricsMap{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}
}

func Run() {
	m := NewMetricsMap()

	tickerCollect := time.NewTicker(pollInterval)
	tickerSend := time.NewTicker(reportInterval)

	defer tickerCollect.Stop()
	defer tickerSend.Stop()
	for {
		select {
		case <-tickerCollect.C:
			m.CollectMetrics()
		case <-tickerSend.C:
			m.SendMetrics("localhost:8080")
		}
	}

}
