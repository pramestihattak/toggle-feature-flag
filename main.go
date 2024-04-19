package main

import (
	"log"
	"net/http"
	"pocfflag/handler"
	"pocfflag/service"
	"pocfflag/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	storage := storage.NewStorage()
	feats, err := storage.GetFeatures()
	if err != nil {
		log.Fatal("crashed")
	}

	fflagHandler := handler.NewFflagHandler(storage, feats)
	service := service.NewService(fflagHandler, storage)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", service.Index)
	http.ListenAndServe(":3000", r)
}
