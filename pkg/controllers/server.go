package controllers

import (
	"log"

	"axis/ecommerce-backend/internal/storage"
)

type Serve struct {
	Db storage.Storage
}

func InitServer() (*Serve, error) {
	db, err := storage.ConnectDb()
	if err != nil {
		return nil, err
	}

	log.Println("Ping Redis")
	if err = storage.Cache.Ping(); err != nil {
		return nil, err
	}

	return &Serve{Db: db}, nil
}
