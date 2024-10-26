package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL    string
	MinioHost      string
	MinioAccessKey string
	MinioSecretKey string
	MinioUseSSL    bool
	ServerAddress  string
}

func NewConfig(environment string) (*Config, error) {
	err := godotenv.Load(environment)
	if err != nil {
		return nil, err
	}

	conf := &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		MinioHost:      os.Getenv("MINIO_HOST"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioUseSSL:    os.Getenv("MINIO_USE_SSL") == "true",
		ServerAddress:  os.Getenv("SERVER_ADDRESS"),
	}
	return conf, nil
}
