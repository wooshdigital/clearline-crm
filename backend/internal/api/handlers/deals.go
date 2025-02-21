package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// DealHandler handles deal pipeline operations.
type DealHandler struct {
	db *sql.DB
}

func NewDealHandler(db *sql.DB) *DealHandler {
	return &DealHandler{db: db}
}

type Deal struct {
	ID          string     `json:"id"`
	ContactID   *string    `json:"contact_id"`
	OwnerID     *string    `json:"owner_id"`
	Title       string     `json:"title"`
	Value       float64    `json:"value"`
	Stage       string     `json:"stage"`
	CloseDate   *string    `json:"close_date"`
	Notes       *string    `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ContactName *string    `json:"contact_name,omitempty"`
}

// List returns all deals with optional stage filter.
func (h *DealHandler) List(w http.ResponseWriter, r *http.Request) {
	stageFilter := r.URL.Query().Get("stage")

	query := `SELECT d.id, d.contact_id, d.owner_id, d.title, d.value, d.stage,
			         d.close_date::text, d.notes, d.created_at, d.updated_at,
			         CONCAT(c.first_name, ' ', c.last_name)
			  FROM deals d
			  LEFT JOIN contacts c ON c.id = d.contact_id
			  WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if stageFilter != "" {
		query += ` AND d.stage = $` + itoa(argIdx)
		args = append(args, stageFilter)
		argIdx++
	}
	query += ` ORDER BY d.created_at DESC`

	rows, err := h.db.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed")
		return
	}
	defer rows.Close()

	deals := []Deal{}
	for rows.Next() {
		var d Deal
		if err := rows.Scan(&d.ID, &d.ContactID, &d.OwnerID, &d.Title, &d.Value, &d.Stage,
			&d.CloseDate, &d.Notes, &d.CreatedAt, &d.UpdatedAt, &d.ContactName); err != nil {
			writeError(w, http.StatusInternalServerError, "scan failed")
			return
		}
		deals = append(deals, d)
	}
	writeJSON(w, http.StatusOK, deals)
}
