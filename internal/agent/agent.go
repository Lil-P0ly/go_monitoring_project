package agent

import (
	"time"

	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/model"
)

type Agent struct {
	cfg     Config
	metrics model.Metrics
}

func New(cfg Config, m model.Metrics) *Agent {
	return &Agent{cfg: cfg, metrics: m}
}

func (a *Agent) Run() {
	poll := time.NewTicker(a.cfg.PollInterval)
	report := time.NewTicker(a.cfg.ReportInterval)
	defer poll.Stop()
	defer report.Stop()
	for {
		select {
		case <-poll.C:
			a.metrics.CollectMetrics()
		case <-report.C:
			a.metrics.SendMetrics(a.cfg.Address)
		}
	}
}