package main

import (
	"log"
	"net/http"

	"github.com/Lil-P0ly/go_monitoring_project/internal/server/handler"
	models "github.com/Lil-P0ly/go_monitoring_project/internal/server/model"
	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	memStorage := models.NewMemStorage()

	msh := handler.NewMSHandlerWithStorage(memStorage)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	})

	r.Get("/", msh.PrintMetricsHTML)
	r.Post("/update/{metrics_type}/{metrics_name}/{metrics_value}", msh.AddValue)
	r.Get("/value/{metrics_type}/{metrics_name}", msh.PrintLastValue)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
