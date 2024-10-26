package services

import (
	"context"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/pkg/minio"
	"github.com/Ayano2000/push/internal/pkg/pgsql"
)

type Services struct {
	Config *config.Config
	DB     *pgsql.Pgsql
	Minio  *minio.Minio
}

func NewServices(ctx context.Context, config *config.Config) (*Services, error) {
	pgsqlClient, err := pgsql.NewPgsql(ctx, config)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.NewMinio(config)
	if err != nil {
		return nil, err
	}

	return &Services{
		Config: config,
		DB:     pgsqlClient,
		Minio:  minioClient,
	}, nil
}

func (s *Services) Cleanup() {
	s.DB.Close()
}
