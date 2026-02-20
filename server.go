package main

import (
	"github.com/h4r5h1l/Chirpy/api"
	"net/http"
)

func main() {
	// Create a new ServeMux to handle incoming HTTP requests
	mux := http.NewServeMux()
	// Create a new HTTP server with the specified address and handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// Create a file server to serve files from the current directory, stripping the "/app" prefix from the URL path
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	// Create a new apiConfig instance to manage the file server hit metrics
	apiCfg := &api.ApiConfig{}
	// Register the file server handler with the middleware to increment the hit count
	mux.Handle("/app/", apiCfg.MiddlewareMetricInc(fileServer))
	// Register handlers for the metrics and reset endpoints
	mux.HandleFunc("GET /admin/metrics", apiCfg.GetFileServerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetFileServerHits)
	// Register a handler for the health check endpoint
	mux.HandleFunc("GET /healthz", readyCheckHandler)
	// Start the HTTP server and listen for incoming requests
	server.ListenAndServe()
}

// readyCheckHandler is a simple HTTP handler that responds with "OK" to indicate that the server is ready to handle requests
func readyCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to indicate that the response is plain text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Write the HTTP status code and response body to indicate that the server is healthy
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
