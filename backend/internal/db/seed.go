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
	}

	// Seed contacts
	type seedContact struct {
		First, Last, Email, Phone, Company, Title, Status string
	}
	contacts := []seedContact{
		{"Alice", "Nguyen", "alice@acmecorp.com", "+1-555-0101", "Acme Corp", "VP of Engineering", "customer"},
		{"Marcus", "Webb", "marcus@globex.io", "+1-555-0202", "Globex", "CTO", "prospect"},
		{"Sarah", "Kim", "sarah@initech.co", "+1-555-0303", "Initech", "Head of IT", "lead"},
		{"David", "Torres", "dtorres@veridian.com", "+1-555-0404", "Veridian Dynamics", "CEO", "prospect"},
		{"Priya", "Patel", "priya@nanoteck.io", "+1-555-0505", "Nanoteck", "Founder", "customer"},
		{"James", "Okonkwo", "james@bluewave.io", "+1-555-0606", "Bluewave", "Director of Ops", "lead"},
	}

	repID := userIDs["rep@clearline.local"]
	contactIDs := []string{}
	for _, c := range contacts {
		var id string
		err := db.QueryRow(
			`INSERT INTO contacts (owner_id, first_name, last_name, email, phone, company, title, status)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
			repID, c.First, c.Last, c.Email, c.Phone, c.Company, c.Title, c.Status,
		).Scan(&id)
		if err != nil {
			return err
		}
		contactIDs = append(contactIDs, id)
	}

	// Seed deals
	type seedDeal struct {
		ContactIdx int
		Title      string
		Value      float64
		Stage      string
		CloseDate  string
	}
	deals := []seedDeal{
		{0, "Acme Enterprise License", 48000, "closed_won", "2024-03-15"},
		{1, "Globex Platform Upgrade", 22500, "proposal", "2024-09-01"},
		{2, "Initech Starter Pack", 5000, "qualification", "2024-10-15"},
		{3, "Veridian Annual Subscription", 75000, "negotiation", "2024-08-30"},
		{4, "Nanoteck Renewal", 18000, "closed_won", "2024-04-01"},
		{5, "Bluewave Discovery", 3000, "prospecting", "2024-11-30"},
	}

	mgrID := userIDs["manager@clearline.local"]
	dealIDs := []string{}
	for _, d := range deals {
		var id string
		err := db.QueryRow(
			`INSERT INTO deals (contact_id, owner_id, title, value, stage, close_date)
			 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
			contactIDs[d.ContactIdx], mgrID, d.Title, d.Value, d.Stage, d.CloseDate,
		).Scan(&id)
		if err != nil {
			return err
		}
		dealIDs = append(dealIDs, id)
	}

	// Seed activities
	type seedActivity struct {
		ContactIdx int
		DealIdx    int
		Type       string
		Subject    string
		Body       string
	}
	activities := []seedActivity{
		{0, 0, "call", "Discovery call with Alice", "Discussed enterprise needs and pricing tier."},
		{0, 0, "email", "Sent proposal to Alice", "Attached the enterprise license proposal PDF."},
		{1, 1, "meeting", "Globex demo session", "Walked through the platform upgrade roadmap."},
		{2, 2, "note", "Initial research", "Initech is evaluating us against two competitors."},
		{3, 3, "email", "Veridian follow-up", "Confirmed interest in annual pricing."},
		{4, 4, "call", "Nanoteck renewal check-in", "Priya wants to add 5 more seats."},
	}

	adminID := userIDs["admin@clearline.local"]
	for _, a := range activities {
		_, err := db.Exec(
			`INSERT INTO activities (user_id, contact_id, deal_id, type, subject, body)
			 VALUES ($1,$2,$3,$4,$5,$6)`,
			adminID, contactIDs[a.ContactIdx], dealIDs[a.DealIdx], a.Type, a.Subject, a.Body,
		)
		if err != nil {
			return err
		}
	}

	log.Println("[seed] inserted demo users, contacts, deals, and activities")
	return nil
}

