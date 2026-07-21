package main

import (
	"fmt"

	"github.com/Lil-P0ly/go_monitoring_project/internal/agent"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/config"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/logger"
	"github.com/Lil-P0ly/go_monitoring_project/internal/agent/model"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := logger.InitLogger(cfg.Level); err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()
	logger.Info("Starting metrics agent",
		zap.String("URL", cfg.Address),
		zap.String("LogLevel", cfg.Level),
	)
	m := model.NewMetricsMap()
	agent.New(cfg, m).Run()
}
