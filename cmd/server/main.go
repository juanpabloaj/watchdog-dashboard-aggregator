package main

import (
	"cmp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/juanpabloaj/watchdogdashboard/internal/application/dashboard"
	"github.com/juanpabloaj/watchdogdashboard/internal/infrastructure/dummyjson"
	"github.com/juanpabloaj/watchdogdashboard/internal/interfaces/httphandler"
)

var (
	buildDate string
	gitHash   string
)

var (
	httpClientTimeoutSeconds = 2
	httpPort                 = "8080"
	dummyjsonURL             = "https://dummyjson.com"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("build date: %s", buildDate)
	log.Printf("git hash: %s", gitHash)

	port := ":" + cmp.Or(os.Getenv("PORT"), httpPort)
	httpClientTimeout := getEnvInt("HTTP_CLIENT_TIMEOUT_SECONDS", httpClientTimeoutSeconds)
	dJSONURL := cmp.Or(os.Getenv("DUMMYJSON_URL"), dummyjsonURL)

	log.Printf("PORT: %s", port)
	log.Printf("HTTP_CLIENT_TIMEOUT_SECONDS: %d", httpClientTimeout)
	log.Printf("DUMMYJSON_URL: %s", dJSONURL)

	httpClient := &http.Client{
		Timeout: time.Duration(httpClientTimeout) * time.Second,
	}

	dummyjsonClient := dummyjson.NewDummyJSONClient(dJSONURL, httpClient)

	dashboardService := dashboard.NewService(dummyjsonClient, dummyjsonClient)

	dashboardHandler := httphandler.NewDashboardHandler(dashboardService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /dashboard/{id}", dashboardHandler.GetDashboard)

	httpServer := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("starting server on port %s", port)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return fallback
}
