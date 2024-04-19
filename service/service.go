package service

import (
	"net/http"
	"pocfflag/handler"
	"pocfflag/storage"
)

type Service struct {
	fflagHandler *handler.FflagHandler
	storage      *storage.Storage
}

func NewService(fflagHandler *handler.FflagHandler, storage *storage.Storage) *Service {
	return &Service{
		fflagHandler: fflagHandler,
		storage:      storage,
	}
}

func (s *Service) Index(w http.ResponseWriter, r *http.Request) {
	fflag := s.fflagHandler.GetFflags()
	if !fflag[handler.FeatureA] {
		w.Write([]byte("disabled"))
		return
	}

	w.Write([]byte("enabled"))
}
