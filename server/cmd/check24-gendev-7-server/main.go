/*
 * CHECK24 GenDev 7 API
 *
 * API for the 7th CHECK24 GenDev challenge providing product offerings from five different internet providers
 *
 * API version: 1.0.0
 */

package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/rotmanjanez/check24-gendev-7/config"
	"github.com/rotmanjanez/check24-gendev-7/internal/api"
	"github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"

	_ "github.com/rotmanjanez/check24-gendev-7/providers/byteme"
	_ "github.com/rotmanjanez/check24-gendev-7/providers/exampleprovider"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	envPath := flag.String("env", ".env", "Path to the environment file")

	flag.Parse()

	if configPath == nil {
		log.Fatal("Config path is nil")
	}
	if envPath == nil {
		log.Fatal("Env path is nil")
	}

	err := godotenv.Load(*envPath)
	if err != nil {
		log.Fatalf("Error loading environment file '%s': %v", *envPath, err)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	providers, err := interfaces.CreateProviders(cfg)
	if err != nil {
		log.Fatalf("Error creating backends: %v", err)
	}

	HealthAPIService := api.NewHealthAPIService()
	HealthAPIController := api.NewHealthAPIController(HealthAPIService)

	SystemAPIService := api.NewSystemAPIService(cfg)
	SystemAPIController := api.NewSystemAPIController(SystemAPIService)

	InternetProductsAPIService := api.NewInternetProductsAPIService(cfg, providers)
	InternetProductsAPIController := api.NewInternetProductsAPIController(InternetProductsAPIService)

	router := api.NewRouter(HealthAPIController, SystemAPIController, InternetProductsAPIController)

	slog.Debug("Using config file", "path", *configPath)
	slog.Debug("Using config", "config", cfg)
	slog.Debug("Using config backends", "backends", cfg.Backends)
	log.Printf("Starting server on %s", cfg.GetAddress())

	log.Fatal(http.ListenAndServe(cfg.GetAddress(), router))
}
