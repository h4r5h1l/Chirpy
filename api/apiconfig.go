package api

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/h4r5h1l/Chirpy/internal/database"
)

// ApiConfig is a struct that holds the hit count for the file server and provides methods to manage it
type ApiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
}

// MiddlewareMetricInc is a HOF that increments the file server hit count and then calls the next handler in the chain
func (cfg *ApiConfig) MiddlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// GetFileServerHits is an HTTP handler that responds with the current hit count for the file server
func (cfg *ApiConfig) GetFileServerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!doctype html>
					<html>
					<body>
						<h1>Welcome, Chirpy Admin</h1>
						<p>Chirpy has been visited %d times!</p>
					</body>
					</html>
					`, cfg.fileserverHits.Load())
}

// ResetFileServerHits resets the hit count and returns an HTML page showing the new value.
func (cfg *ApiConfig) ResetFileServerHits(w http.ResponseWriter, r *http.Request) {
	// Only allow this dangerous operation in development
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Forbidden")
		return
	}

	// Delete all users from the database (do not touch schema)
	if cfg.DB != nil {
		if err := cfg.DB.DeleteAllUsers(r.Context()); err != nil {
			http.Error(w, "failed to delete users", http.StatusInternalServerError)
			return
		}
	}

	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!doctype html>
					<html>
					<body>
						<h1>Welcome, Chirpy Admin</h1>
						<p>All users deleted and visits reset to %d.</p>
					</body>
					</html>
					`, cfg.fileserverHits.Load())
}
