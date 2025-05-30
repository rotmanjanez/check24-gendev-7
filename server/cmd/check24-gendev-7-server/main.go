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
	"github.com/rotmanjanez/check24-gendev-7/pkg/cache"
	"github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	"github.com/rotmanjanez/check24-gendev-7/pkg/logger"
	"github.com/rotmanjanez/check24-gendev-7/pkg/provider"

	_ "github.com/rotmanjanez/check24-gendev-7/providers/byteme"
	_ "github.com/rotmanjanez/check24-gendev-7/providers/exampleprovider"
	_ "github.com/rotmanjanez/check24-gendev-7/providers/pingperfect"
	_ "github.com/rotmanjanez/check24-gendev-7/providers/servusspeed"
	_ "github.com/rotmanjanez/check24-gendev-7/providers/verbyndich"
	_ "github.com/rotmanjanez/check24-gendev-7/providers/webwunder"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	envPath := flag.String("env", ".env", "Path to the environment file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	localdev := flag.Bool("localdev", false, "Enable local development mode")
	jsonLogger := flag.Bool("json-logger", false, "Use JSON logger format")

	flag.Parse()

	if configPath == nil {
		log.Fatal("Config path is nil")
	}
	if envPath == nil {
		log.Fatal("Env path is nil")
	}
	if debug == nil {
		log.Fatal("Debug flag is nil")
	}
	if localdev == nil {
		log.Fatal("Localdev flag is nil")
	}
	if jsonLogger == nil {
		log.Fatal("JSON logger flag is nil")
	}

	var handler slog.Handler
	logOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if *debug {
		logOptions.Level = slog.LevelDebug
	}

	if *jsonLogger {
		handler = slog.NewJSONHandler(log.Writer(), logOptions)
	} else {
		// Use our custom text handler
		handler = logger.NewTextHandler(log.Writer(), logOptions)
	}

	slog.SetDefault(slog.New(handler))

	err := godotenv.Load(*envPath)
	if err != nil {
		log.Fatalf("Error loading environment file '%s': %v", *envPath, err)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if cfg == nil {
		log.Fatal("Config is nil")
	}

	if cfg.GetAddress() == "" {
		log.Fatal("Config address is empty")
	}

	if len(cfg.Backends) == 0 {
		slog.Warn("No backends configured")
	}

	if !cfg.UseInProcessCache && *localdev {
		slog.Info("Overriding UseInProcessCache to true in local development mode")
		cfg.UseInProcessCache = true
	}

	var cacheFactory interfaces.CacheFactory

	if cfg.UseInProcessCache {
		cacheFactory = cache.NewInstanceCacheFactory()
	} else {
		if cfg.Redis == nil {
			log.Fatal("Redis configuration is missing in the config file")
		}
		cacheFactory, err = cache.NewRedisCacheFactory(cfg.Redis)
		if err != nil {
			log.Fatalf("Failed to create Redis cache factory: %v", err)
		}
	}

	providers, err := provider.CreateProviders(cacheFactory, cfg)
	if err != nil {
		log.Fatalf("Error creating backends: %v", err)
	}

	HealthAPIService := api.NewHealthAPIService()
	HealthAPIController := api.NewHealthAPIController(HealthAPIService)

	SystemAPIService := api.NewSystemAPIService(cfg)
	SystemAPIController := api.NewSystemAPIController(SystemAPIService)

	cache, err := cacheFactory.Create("check24-gendev-7")
	if err != nil {
		log.Fatalf("Error creating cache: %v", err)
	}
	queue, err := cacheFactory.Create("check24-gendev-7-queue")
	if err != nil {
		log.Fatalf("Error creating queue: %v", err)
	}
	InternetProductsAPIService := api.NewInternetProductsAPIService(cfg, cache, queue, providers)
	InternetProductsAPIController := api.NewInternetProductsAPIController(InternetProductsAPIService)

	router := api.NewRouter(HealthAPIController, SystemAPIController, InternetProductsAPIController)

	slog.Debug("Using config file", "path", *configPath)
	slog.Debug("Using config backends", "backends", cfg.Backends)
	log.Printf("Starting server on %s", cfg.GetAddress())

	log.Fatal(http.ListenAndServe(cfg.GetAddress(), router))
}
