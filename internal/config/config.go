package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL    string
	LogFilePath    string
	MinioHost      string
	MinioAccessKey string
	MinioSecretKey string
	MinioUseSSL    bool
	ServerAddress  string
}

func NewConfig(env string) (*Config, error) {
	err := godotenv.Load(env)
	if err != nil {
		return nil, err
	}

	conf := &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		LogFilePath:    os.Getenv("LOG_FILE_PATH"),
		MinioHost:      os.Getenv("MINIO_HOST"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioUseSSL:    os.Getenv("MINIO_USE_SSL") == "true",
		ServerAddress:  os.Getenv("SERVER_ADDRESS"),
	}
	return conf, nil
}
