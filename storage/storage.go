package storage

import (
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

var (
	getFeatureSQL = `
		SELECT
			feature,
			is_enabled
		FROM toggle_features;
		`
)

func NewListener() (*pq.Listener, error) {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println("err", err.Error())
		}
	}

	listener := pq.NewListener("port=5438 user=postgres dbname=postgres password=secret sslmode=disable", 10*time.Second, time.Minute, reportProblem)
	return listener, nil
}

func NewStorage() *Storage {
	db, err := sqlx.Connect("postgres", "port=5438 user=postgres dbname=postgres password=secret sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	return &Storage{
		db: db,
	}

}

func (s *Storage) GetFeatures() ([]*StorageFeature, error) {
	rows, err := s.db.Query(getFeatureSQL)
	if err != nil {
		return nil, errors.New("fail to query")
	}

	defer rows.Close()

	feats := []*StorageFeature{}
	for rows.Next() {
		var feat StorageFeature
		if err := rows.Scan(
			&feat.Feature,
			&feat.IsEnabled,
		); err != nil {
			return nil, errors.New("fail to query")
		}

		feats = append(feats, &feat)
	}

	return feats, nil
}
