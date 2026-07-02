package agent

import (
	"time"

	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/model"
)

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second

func Run(m model.Metrics) {

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
