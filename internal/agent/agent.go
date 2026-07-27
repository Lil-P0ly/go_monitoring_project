package agent

import (
	"time"

	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/config"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/model"
)

type Agent struct {
	cfg     config.Config
	metrics model.MetricSender
}

func New(cfg config.Config, m model.MetricSender) *Agent {
	return &Agent{cfg: cfg, metrics: m}
}

func (a *Agent) Run() {
	poll := time.NewTicker(time.Duration(a.cfg.PollInterval) * time.Second)
	report := time.NewTicker(time.Duration(a.cfg.ReportInterval) * time.Second)
	defer poll.Stop()
	defer report.Stop()
	for {
		select {
		case <-poll.C:
			a.metrics.CollectMetrics()
		case <-report.C:
			a.metrics.SendMetricJSON(a.cfg.Address)
		}
	}
}
