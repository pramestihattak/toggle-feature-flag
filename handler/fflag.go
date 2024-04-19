package handler

import (
	"encoding/json"
	"log"
	"pocfflag/storage"

	"github.com/lib/pq"
)

type FflagHandler struct {
	Fflags  map[string]bool
	storage *storage.Storage
}

const (
	FeatureA = "FEATURE_A"
	FeatureB = "FEATURE_B"
)

func NewFflagHandler(s *storage.Storage, feats []*storage.StorageFeature) *FflagHandler {
	d := map[string]bool{}

	for _, f := range feats {
		d[f.Feature] = f.IsEnabled
	}

	h := &FflagHandler{Fflags: d, storage: s}

	listener, err := storage.NewListener()
	if err != nil {
		log.Println("fail to initialize listener")
	}

	if err := listener.Listen("toggle_feature"); err != nil {
		log.Println("fail to initialize listener")
	}

	feat := make(chan storage.StorageFeature)
	go h.CheckToggleFeatureChange(listener, feat)

	go func() {
		for f := range feat {
			log.Println("feature changes", f)
			h.UpdateToggleFeatures(f)
		}
	}()

	return h
}

func (h *FflagHandler) CheckToggleFeatureChange(l *pq.Listener, feat chan storage.StorageFeature) {
	var featData struct {
		Action string `json:"action"`
		Table  string `json:"table"`
		Data   struct {
			Feature   string `json:"feature"`
			IsEnabled bool   `json:"is_enabled"`
		} `json:"data"`
	}
	for n := range l.Notify {
		if n == nil {
			continue
		}

		var featDataMap map[string]interface{}
		if err := json.Unmarshal([]byte(n.Extra), &featDataMap); err != nil {
			continue
		}

		featDataStr, err := json.Marshal(featDataMap)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(featDataStr, &featData); err != nil {
			continue
		}

		sFeat := storage.StorageFeature{
			Feature:   featData.Data.Feature,
			IsEnabled: featData.Data.IsEnabled,
		}

		feat <- sFeat
	}
}

func (h *FflagHandler) UpdateToggleFeatures(feat storage.StorageFeature) {
	feats := h.GetFflags()
	for k := range feats {
		if k == feat.Feature {
			feats[k] = feat.IsEnabled
		}
	}

	h.Fflags = feats

	log.Println("DONE UPDATING")
}

func (h *FflagHandler) GetFflags() map[string]bool {
	return h.Fflags
}
