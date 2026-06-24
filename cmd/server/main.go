package main

import (
	"net/http"

	"github.com/Lil-P0ly/go_monitoring_project/internal/server/handler"
)

func main() {
	mux := http.NewServeMux()

	msh := handler.NewMSHandler()

	mux.HandleFunc("/update/{metrics_type}/{metric_name}/{metrics_value}", msh.AddValue)

	mux.HandleFunc("/update/{metrics_type}/", msh.NotFound)

	mux.HandleFunc("/metrics", msh.PrintMetrics)

	http.ListenAndServe(":8080", mux)

}
