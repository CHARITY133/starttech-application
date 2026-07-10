// Example backend entrypoint satisfying the assessment's requirements.
// Replace with your actual application logic from the much-to-do repo,
// keeping the /api/v1/health endpoint, port 8080, and env var contract.
package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	redisHost := os.Getenv("REDIS_HOST")
	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if redisHost == "" {
		slog.Warn("REDIS_HOST not set")
	}
	if mongoURI == "" {
		slog.Warn("MONGO_URI not set")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// TODO: mount real application routes here, e.g.:
	// mux.HandleFunc("/api/v1/tasks", tasksHandler)

	slog.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server exited", "error", err)
		os.Exit(1)
	}
}
