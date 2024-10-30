package main

import (
	"fmt"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/pkg/logger"
	"github.com/Ayano2000/push/pkg/router"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument 'environment'. Usage: make run <development|production>")
		os.Exit(1)
	}

	conf, err := config.NewConfig(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load Config: %v\n", err)
		os.Exit(1)
	}

	logger.MustSetupLogger(conf)

	handler, err := handlers.NewHandler(conf)
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
