package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	Name     string
	Email    string
	Password string
	Role     string
}

// Seed inserts demo data if the users table is empty.
func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		log.Println("[seed] data already present, skipping")
		return nil
	}

	users := []seedUser{
		{"Admin User", "admin@clearline.local", "admin123", "admin"},
		{"Jane Manager", "manager@clearline.local", "manager123", "manager"},
		{"Bob Rep", "rep@clearline.local", "rep123", "rep"},
	}

	userIDs := map[string]string{}
	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		var id string
		err = db.QueryRow(
			`INSERT INTO users (name, email, password_hash, role) VALUES ($1,$2,$3,$4) RETURNING id`,
			u.Name, u.Email, string(hash), u.Role,
		).Scan(&id)
		if err != nil {
			return err
		}
		userIDs[u.Email] = id