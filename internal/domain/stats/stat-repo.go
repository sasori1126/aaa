package stats

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
	"encoding/json"
	"time"
)

type StatRepo interface {
	GetTotalStat(limit int) (*models.Stat, error)
}

type DefaultStatDbRepo struct {
	db storage.Storage
}

func (d DefaultStatDbRepo) GetTotalStat(limit int) (*models.Stat, error) {
	data, err := storage.Cache.RememberWithTime("dashboard", func(f interface{}) (string, error) {
		stat, err := d.db.GetStats(f.(int))
		if err != nil {
			return "", err
		}
		marshalledData, err := json.Marshal(stat)
		if err != nil {
			return "", err
		}
		return string(marshalledData), nil
	}, limit, 3*time.Hour)
	if err != nil {
		return nil, err
	}
	stat := &models.Stat{}
	err = json.Unmarshal([]byte(data), stat)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func NewDefaultStatDbRepo(db storage.Storage) StatRepo {
	return &DefaultStatDbRepo{db: db}
}
