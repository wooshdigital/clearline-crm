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

// Get returns a single deal by ID.
func (h *DealHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var d Deal
	err := h.db.QueryRow(
		`SELECT d.id, d.contact_id, d.owner_id, d.title, d.value, d.stage,
		        d.close_date::text, d.notes, d.created_at, d.updated_at,
		        CONCAT(c.first_name, ' ', c.last_name)
		 FROM deals d LEFT JOIN contacts c ON c.id = d.contact_id WHERE d.id = $1`, id,
	).Scan(&d.ID, &d.ContactID, &d.OwnerID, &d.Title, &d.Value, &d.Stage,
		&d.CloseDate, &d.Notes, &d.CreatedAt, &d.UpdatedAt, &d.ContactName)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "deal not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed")
		return
	}
	writeJSON(w, http.StatusOK, d)
}

// Create adds a new deal.
func (h *DealHandler) Create(w http.ResponseWriter, r *http.Request) {
	var d Deal
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if strings.TrimSpace(d.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if d.Stage == "" {
		d.Stage = "prospecting"
	}

	var id string
	err := h.db.QueryRow(
		`INSERT INTO deals (contact_id, owner_id, title, value, stage, close_date, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		d.ContactID, d.OwnerID, d.Title, d.Value, d.Stage, d.CloseDate, d.Notes,
	).Scan(&id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "insert failed")
		return
	}
	d.ID = id
	writeJSON(w, http.StatusCreated, d)
}

// Update modifies an existing deal.
func (h *DealHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var d Deal
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	_, err := h.db.Exec(
		`UPDATE deals SET title=$1, value=$2, stage=$3, close_date=$4, notes=$5, updated_at=NOW() WHERE id=$6`,
		d.Title, d.Value, d.Stage, d.CloseDate, d.Notes, id,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update failed")
		return
	}
	d.ID = id
	writeJSON(w, http.StatusOK, d)
}

