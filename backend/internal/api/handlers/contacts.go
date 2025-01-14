package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// ContactHandler handles contact CRUD operations.
type ContactHandler struct {
	db *sql.DB
}

func NewContactHandler(db *sql.DB) *ContactHandler {
	return &ContactHandler{db: db}
}

type Contact struct {
	ID        string     `json:"id"`
	OwnerID   *string    `json:"owner_id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     *string    `json:"email"`
	Phone     *string    `json:"phone"`
	Company   *string    `json:"company"`
	Title     *string    `json:"title"`
	Status    string     `json:"status"`
	Notes     *string    `json:"notes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// List returns all contacts, optionally filtered by status or search query.
func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	statusFilter := q.Get("status")
	searchQuery := q.Get("q")

	baseQuery := `SELECT id, owner_id, first_name, last_name, email, phone, company, title, status, notes, created_at, updated_at
				  FROM contacts WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if statusFilter != "" {
		baseQuery += ` AND status = $` + itoa(argIdx)
		args = append(args, statusFilter)
		argIdx++
	}

	if searchQuery != "" {
		like := "%" + strings.ToLower(searchQuery) + "%"
		baseQuery += ` AND (LOWER(first_name) LIKE $` + itoa(argIdx) +
			` OR LOWER(last_name) LIKE $` + itoa(argIdx) +
			` OR LOWER(email) LIKE $` + itoa(argIdx) +
			` OR LOWER(company) LIKE $` + itoa(argIdx) + `)`
		args = append(args, like)
		argIdx++
	}

	baseQuery += ` ORDER BY created_at DESC`

	rows, err := h.db.Query(baseQuery, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed")
		return
	}
	defer rows.Close()

	contacts := []Contact{}
	for rows.Next() {
		var c Contact
		if err := rows.Scan(&c.ID, &c.OwnerID, &c.FirstName, &c.LastName, &c.Email, &c.Phone,
			&c.Company, &c.Title, &c.Status, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "scan failed")
			return
		}
		contacts = append(contacts, c)
	}
	writeJSON(w, http.StatusOK, contacts)
}

// Get returns a single contact by ID.
func (h *ContactHandler) Get(w http.ResponseWriter, r *http.Request) {