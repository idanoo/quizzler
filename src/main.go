package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"quizzler/cache"
	"quizzler/database"
	"quizzler/handlers"
	"quizzler/middleware"
)

func main() {
	// Setup logging
	logLevel := slog.LevelInfo
	if os.Getenv("DEBUG") == "true" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Initialize database
	if err := database.Init(); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize cache
	if err := cache.Init(); err != nil {
		slog.Warn("Failed to connect to Redis, caching disabled", "error", err)
	} else {
		defer cache.Close()
	}

	// API routes
	api := http.NewServeMux()
	api.Handle("POST /register", http.HandlerFunc(handlers.Register))
	api.Handle("POST /login", http.HandlerFunc(handlers.Login))

	api.Handle("GET /decks", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetDecks)))
	api.Handle("GET /decks/public", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPublicDecks)))
	api.Handle("POST /decks", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateDeck)))
	api.Handle("GET /decks/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetDeck)))
	api.Handle("PUT /decks/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateDeck)))
	api.Handle("DELETE /decks/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteDeck)))

	api.Handle("GET /decks/{deckId}/cards", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetCards)))
	api.Handle("POST /decks/{deckId}/cards", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateCard)))
	api.Handle("POST /decks/{deckId}/cards/import", middleware.AuthMiddleware(http.HandlerFunc(handlers.ImportCards)))
	api.Handle("GET /cards/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetCard)))
	api.Handle("PUT /cards/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateCard)))
	api.Handle("DELETE /cards/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteCard)))

	// Main router
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", middleware.JSONMiddleware(api)))
	mux.Handle("/", spaHandler("/web"))

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8132"
	}

	slog.Info("Server starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}
}

// spaHandler serves static files and falls back to index.html for SPA routes
func spaHandler(staticPath string) http.Handler {
	fs := http.Dir(staticPath)
	fileServer := http.FileServer(fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Check if file exists
		if f, err := fs.Open(path); err == nil {
			f.Close()
			// Check if it's a file (has extension) or exact match
			if strings.Contains(filepath.Base(path), ".") || path == "/" {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		// Serve index.html for SPA routes
		http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
	})
}
