package main

import (
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/model"
)

func main() {
	m := model.NewMetricsMap()

	agent.Run(m)
}
