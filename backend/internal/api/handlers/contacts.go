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
	id := r.PathValue("id")
	var c Contact
	err := h.db.QueryRow(
		`SELECT id, owner_id, first_name, last_name, email, phone, company, title, status, notes, created_at, updated_at
		 FROM contacts WHERE id = $1`, id,
	).Scan(&c.ID, &c.OwnerID, &c.FirstName, &c.LastName, &c.Email, &c.Phone,
		&c.Company, &c.Title, &c.Status, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "contact not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// Create adds a new contact.
func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if strings.TrimSpace(c.FirstName) == "" || strings.TrimSpace(c.LastName) == "" {
		writeError(w, http.StatusBadRequest, "first_name and last_name are required")
		return
	}
	if c.Status == "" {
		c.Status = "lead"
	}

	var id string
	err := h.db.QueryRow(
		`INSERT INTO contacts (owner_id, first_name, last_name, email, phone, company, title, status, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		c.OwnerID, c.FirstName, c.LastName, c.Email, c.Phone, c.Company, c.Title, c.Status, c.Notes,
	).Scan(&id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "insert failed")
		return
	}
	c.ID = id
	writeJSON(w, http.StatusCreated, c)
}

// Update modifies an existing contact.
func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	_, err := h.db.Exec(
		`UPDATE contacts SET first_name=$1, last_name=$2, email=$3, phone=$4, company=$5,
		 title=$6, status=$7, notes=$8, updated_at=NOW() WHERE id=$9`,
		c.FirstName, c.LastName, c.Email, c.Phone, c.Company, c.Title, c.Status, c.Notes, id,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update failed")
		return
	}
	c.ID = id
	writeJSON(w, http.StatusOK, c)
}

// Delete removes a contact.
func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err := h.db.Exec(`DELETE FROM contacts WHERE id = $1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "delete failed")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "contact not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
