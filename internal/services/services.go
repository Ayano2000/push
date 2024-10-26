package services

import (
	"context"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/pkg/minio"
	"github.com/jackc/pgx/v5"
)

type Services struct {
	Config *config.Config
	DB     *pgx.Conn
	Minio  *minio.Minio
}

func NewServices(config *config.Config) (*Services, error) {
	pgsqlConn, err := pgx.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.NewMinio(config)
	if err != nil {
		return nil, err
	}

	return &Services{
		Config: config,
		DB:     pgsqlConn,
		Minio:  minioClient,
	}, nil
}

func (s *Services) Cleanup() error {
	return s.DB.Close(context.Background())
}
