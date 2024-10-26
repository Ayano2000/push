package handlers

import (
	"context"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/services"
)

type Handler struct {
	Config   *config.Config
	Services *services.Services
}

func NewHandler(ctx context.Context, config *config.Config) (*Handler, error) {
	s, err := services.NewServices(ctx, config)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		Config:   config,
		Services: s,
	}
	return handler, nil
}
