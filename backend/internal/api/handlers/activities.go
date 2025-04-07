package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// ActivityHandler handles activity timeline operations.
type ActivityHandler struct {
	db *sql.DB
}

func NewActivityHandler(db *sql.DB) *ActivityHandler {
	return &ActivityHandler{db: db}
}

type Activity struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"user_id"`
	ContactID   *string   `json:"contact_id"`
	DealID      *string   `json:"deal_id"`
	Type        string    `json:"type"`
	Subject     string    `json:"subject"`
	Body        *string   `json:"body"`
	OccurredAt  time.Time `json:"occurred_at"`
	CreatedAt   time.Time `json:"created_at"`
	UserName    *string   `json:"user_name,omitempty"`
}

// List returns activities with optional filters.
func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	contactID := q.Get("contact_id")
	dealID := q.Get("deal_id")
	actType := q.Get("type")

	query := `SELECT a.id, a.user_id, a.contact_id, a.deal_id, a.type, a.subject, a.body,
			         a.occurred_at, a.created_at, u.name
			  FROM activities a
			  LEFT JOIN users u ON u.id = a.user_id
			  WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if contactID != "" {
		query += ` AND a.contact_id = $` + itoa(argIdx)
		args = append(args, contactID)
		argIdx++
	}
	if dealID != "" {
		query += ` AND a.deal_id = $` + itoa(argIdx)
		args = append(args, dealID)
		argIdx++
	}
	if actType != "" {
		query += ` AND a.type = $` + itoa(argIdx)
		args = append(args, actType)
		argIdx++
	}
	query += ` ORDER BY a.occurred_at DESC LIMIT 100`

	rows, err := h.db.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed")
		return
	}
	defer rows.Close()

	activities := []Activity{}
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.UserID, &a.ContactID, &a.DealID, &a.Type,
			&a.Subject, &a.Body, &a.OccurredAt, &a.CreatedAt, &a.UserName); err != nil {
			writeError(w, http.StatusInternalServerError, "scan failed")
			return
		}
		activities = append(activities, a)
	}
	writeJSON(w, http.StatusOK, activities)
}

// Create logs a new activity.
func (h *ActivityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var a Activity
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if a.Type == "" || a.Subject == "" {
		writeError(w, http.StatusBadRequest, "type and subject are required")
		return
	}

	var id string
	err := h.db.QueryRow(
		`INSERT INTO activities (user_id, contact_id, deal_id, type, subject, body)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		a.UserID, a.ContactID, a.DealID, a.Type, a.Subject, a.Body,
	).Scan(&id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "insert failed")
		return
	}
	a.ID = id
	writeJSON(w, http.StatusCreated, a)
}
