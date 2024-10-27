package main

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/internal/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"net/http"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	if len(os.Args) < 2 {
		fmt.Println("Missing argument 'environment'. Usage: make run <development|production>")
		return
	}

	conf, err := config.NewConfig(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load Config: %v\n", err)
		os.Exit(1)
	}

	handler, err := handlers.NewHandler(context.Background(), conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Handler: %v\n", err)
		os.Exit(1)
	}

	defer handler.Services.Cleanup()

	dmux, err := router.RegisterRoutes(handler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to register router: %v\n", err)
		os.Exit(1)
	}

	server := http.Server{
		Addr:    conf.ServerAddress,
		Handler: dmux,
	}

	fmt.Fprintf(os.Stdout, "Server is running on: %s", conf.ServerAddress)
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stdout, "Failed to start server: %v", err)
	}
}
