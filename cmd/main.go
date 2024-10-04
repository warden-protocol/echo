package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	log "github.com/rs/zerolog"
)

type EndpointResult struct {
	StatusCode int    `json:"status_code"`
	Reachable  bool   `json:"reachable"`
	Error      string `json:"error,omitempty"`
}

type HealthStatus struct {
	Status           string                    `json:"status"`
	Endpoints        map[string]EndpointResult `json:"endpoints"`
	AdditionalChecks map[string]bool           `json:"additional_checks"`
}

type HealthCheck struct {
	log    log.Logger
	config Config
}

const (
	readTimeout  = 5
	writeTimeout = 10
	idleTimeout  = 15
)

func main() {
	log := log.New(
		log.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(log.InfoLevel).With().Timestamp().Logger()

	cfg, err := LoadConfig()
	if err != nil {
		log.Error().Msgf("Failed to load config: %v", err)
	}

	h := HealthCheck{log: log, config: cfg}

	log.Info().Msg("Starting health check server on port 10010...")
	http.HandleFunc("/health", h.healthHandler)
	server := &http.Server{
		Addr:         ":10010",
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
	}
	if err = server.ListenAndServe(); err != nil {
		// log.Fatalf("Server failed to start: %v", err)
		log.Fatal().Msgf("Server failed to start: %v", err)
	}
}

func (h *HealthCheck) healthHandler(w http.ResponseWriter, _ *http.Request) {
	endpointResults := checkEndpoints(h.config.Endpoints)
	additionalChecks := make(map[string]bool)

	overallStatus := "OK"
	for _, res := range endpointResults {
		if !res.Reachable {
			overallStatus = "DOWN"
			break
		}
	}
	if overallStatus == "OK" {
		additionalChecks = performAdditionalChecks(h.config)
		for _, check := range additionalChecks {
			if check {
				overallStatus = "DOWN"
				break
			}
		}
	}

	healthStatus := HealthStatus{
		Status:           overallStatus,
		Endpoints:        endpointResults,
		AdditionalChecks: additionalChecks,
	}

	statusCode := http.StatusOK
	if overallStatus != "OK" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(healthStatus); err != nil {
		h.log.Error().Msgf("Failed to encode health status: %v", err)
	}
	h.log.Info().Msgf("Health check status: %v", healthStatus)
}
