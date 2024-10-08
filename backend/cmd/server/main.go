package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/yourorg/clearline-crm/internal/api"
	"github.com/yourorg/clearline-crm/internal/db"
)

func main() {
	dsn := getEnv("DATABASE_URL", "postgres://crm:crmsecret@localhost:5432/clearline?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	port := getEnv("PORT", "8080")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:4321")

	// Retry DB connection up to 30 times with 2s delay
	var database *sql.DB
	var err error
	for i := 0; i < 30; i++ {
		database, err = sql.Open("postgres", dsn)
		if err == nil {
			err = database.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("[db] waiting for postgres (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("[db] failed to connect after 30 attempts: %v", err)
	}
	defer database.Close()

	log.Println("[db] connected to PostgreSQL")

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatalf("[db] migration failed: %v", err)
	}
	log.Println("[db] migrations applied")

	// Seed demo data
	if err := db.Seed(database); err != nil {
		log.Printf("[db] seed warning: %v", err)
	}
	log.Println("[db] demo data seeded")

	// Build router
	router := api.NewRouter(database, jwtSecret, frontendURL)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("[server] listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("[server] failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
