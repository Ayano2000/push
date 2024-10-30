package services

import (
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/pkg/storage"
)

type Services struct {
	Config *config.Config
	DB     storage.DatabaseHandler
	Minio  storage.ObjectStoreHandler
}

func NewServices(config *config.Config) (*Services, error) {
	pgsql, err := storage.NewPostgresDB(config)
	if err != nil {
		return nil, err
	}

	minio, err := storage.NewMinIOStorage(config)
	if err != nil {
		return nil, err
	}

	return &Services{
		Config: config,
		DB:     pgsql,
		Minio:  minio,
	}, nil
}

func (s *Services) Cleanup() {
	s.DB.Close()
}
