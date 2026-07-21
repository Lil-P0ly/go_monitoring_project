package main

import (
	"fmt"
	"os"

	"github.com/Lil-P0ly/go_monitoring_project/internal/server"
	"github.com/Lil-P0ly/go_monitoring_project/internal/server/config"
	"github.com/Lil-P0ly/go_monitoring_project/internal/server/handler"
	"github.com/Lil-P0ly/go_monitoring_project/internal/server/logger"
	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = logger.InitLogger(cfg.Level)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer logger.Sync()

	logger.Info("Starting metrics server",
		zap.String("URL", cfg.Address),
		zap.String("LogLevel", cfg.Level),
	)

	storage := models.NewMemStorage()
	h := handler.NewMSHandlerWithStorage(storage)

	srv := server.New(cfg, h)
	if err := srv.Run(); err != nil {
		logger.Fatal("Failed to start server",
			zap.Error(err),
		)
	}
}
