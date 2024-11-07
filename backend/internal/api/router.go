package api

import (
	"database/sql"
	"net/http"

	"github.com/yourorg/clearline-crm/internal/api/handlers"
	"github.com/yourorg/clearline-crm/internal/middleware"
)

// NewRouter wires up all routes and returns an http.Handler.
func NewRouter(db *sql.DB, jwtSecret, frontendURL string) http.Handler {
	mux := http.NewServeMux()

	// Health check — no auth required
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth
	authHandler := handlers.NewAuthHandler(db, jwtSecret)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// Protected routes
	auth := middleware.RequireAuth(jwtSecret)

	// Contacts
	contactHandler := handlers.NewContactHandler(db)
	mux.Handle("GET /api/contacts", auth(http.HandlerFunc(contactHandler.List)))
	mux.Handle("POST /api/contacts", auth(http.HandlerFunc(contactHandler.Create)))
	mux.Handle("GET /api/contacts/{id}", auth(http.HandlerFunc(contactHandler.Get)))
	mux.Handle("PUT /api/contacts/{id}", auth(http.HandlerFunc(contactHandler.Update)))
	mux.Handle("DELETE /api/contacts/{id}", auth(http.HandlerFunc(contactHandler.Delete)))

	// Deals
	dealHandler := handlers.NewDealHandler(db)
	mux.Handle("GET /api/deals", auth(http.HandlerFunc(dealHandler.List)))
	mux.Handle("POST /api/deals", auth(http.HandlerFunc(dealHandler.Create)))
	mux.Handle("GET /api/deals/{id}", auth(http.HandlerFunc(dealHandler.Get)))
	mux.Handle("PUT /api/deals/{id}", auth(http.HandlerFunc(dealHandler.Update)))

	// Activities
	activityHandler := handlers.NewActivityHandler(db)
	mux.Handle("GET /api/activities", auth(http.HandlerFunc(activityHandler.List)))
	mux.Handle("POST /api/activities", auth(http.HandlerFunc(activityHandler.Create)))

	// Users
	userHandler := handlers.NewUserHandler(db)
	mux.Handle("GET /api/users", auth(http.HandlerFunc(userHandler.List)))

	// Wrap with CORS middleware
	return middleware.CORS(frontendURL)(mux)
}
