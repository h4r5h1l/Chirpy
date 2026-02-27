package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/h4r5h1l/Chirpy/api"
	"github.com/h4r5h1l/Chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load() // Load environment variables from .env file
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		panic("DB_URL environment variable is required")
	}
	// Initialize the database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Create a new database queries instance to interact with the database
	dbQueries := database.New(db)
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
	apiCfg := &api.ApiConfig{
		DB:       dbQueries,
		Platform: os.Getenv("PLATFORM"),
	}
	// Register the file server handler with the middleware to increment the hit count
	mux.Handle("/app/", apiCfg.MiddlewareMetricInc(fileServer))
	// Register handlers for the metrics and reset endpoints
	mux.HandleFunc("GET /admin/metrics", apiCfg.GetFileServerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetFileServerHits)
	// Register a handler for the api calls
	mux.HandleFunc("POST /api/users", handlerUsers(dbQueries))
	mux.HandleFunc("POST /api/chirps", handlerChirps(dbQueries))
	mux.HandleFunc("GET /api/chirps", handlerGetChirps(dbQueries))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlerGetChirpById(dbQueries))
	// Register a handler for the health check endpoint
	mux.HandleFunc("GET /healthz", readyCheckHandler)
	// Start the HTTP server and listen for incoming requests
	server.ListenAndServe()
}
