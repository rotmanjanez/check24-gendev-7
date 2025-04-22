/*
 * CHECK24 GenDev 7 API
 *
 * API for the 7th CHECK24 GenDev challenge providing product offerings from five different internet providers
 *
 * API version: 1.0.0
 */

package main

import (
	"log"
	"net/http"

	"github.com/rotmanjanez/check24-gendev-7/config"
	"github.com/rotmanjanez/check24-gendev-7/internal/api"
)

func main() {
	log.Printf("Server started")

	cfg := config.NewConfig()

	HealthAPIService := api.NewHealthAPIService()
	HealthAPIController := api.NewHealthAPIController(HealthAPIService)

	router := api.NewRouter(HealthAPIController)

	log.Fatal(http.ListenAndServe(cfg.GetAddress(), router))
}
