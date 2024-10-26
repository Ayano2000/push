package main

import (
	"fmt"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/internal/routes"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument 'environment'. Usage: make run <development|production>")
		return
	}

	conf, err := config.NewConfig(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load Config: %v\n", err)
		os.Exit(1)
	}

	handler, err := handlers.NewHandler(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Handler: %v\n", err)
		os.Exit(1)
	}

	defer handler.Services.Cleanup()

	mux := http.NewServeMux()
	err = routes.RegisterRoutes(mux, handler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to register routes: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Server is running on: %s", conf.ServerAddress)
	if err := http.ListenAndServe(conf.ServerAddress, mux); err != nil {
		fmt.Fprintf(os.Stdout, "Failed to start server: %v", err)
	}
}
