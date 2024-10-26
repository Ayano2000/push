package services

import (
	"context"
	"github.com/Ayano2000/push/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Services struct {
	Config *config.Config
	DB     *pgx.Conn
	Minio  *minio.Client
}

func NewServices(config *config.Config) (*Services, error) {
	pgsqlConn, err := pgx.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	minioConn, err := minio.New(config.MinioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: config.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Services{
		Config: config,
		DB:     pgsqlConn,
		Minio:  minioConn,
	}, nil
}

func (s *Services) Cleanup() error {
	return s.DB.Close(context.Background())
}
