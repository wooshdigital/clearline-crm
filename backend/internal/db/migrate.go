package db

import "database/sql"

// Migrate runs all DDL statements to create the schema if it doesn't exist.
func Migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'rep' CHECK (role IN ('admin','manager','rep')),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS contacts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			company TEXT,
			title TEXT,
			status TEXT NOT NULL DEFAULT 'lead' CHECK (status IN ('lead','prospect','customer','churned')),
			notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS deals (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			contact_id UUID REFERENCES contacts(id) ON DELETE SET NULL,
			owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
			title TEXT NOT NULL,
			value NUMERIC(12,2) NOT NULL DEFAULT 0,
			stage TEXT NOT NULL DEFAULT 'prospecting' CHECK (stage IN ('prospecting','qualification','proposal','negotiation','closed_won','closed_lost')),
			close_date DATE,
			notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS activities (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE SET NULL,
			contact_id UUID REFERENCES contacts(id) ON DELETE CASCADE,
			deal_id UUID REFERENCES deals(id) ON DELETE CASCADE,
			type TEXT NOT NULL CHECK (type IN ('call','email','meeting','note','task')),
			subject TEXT NOT NULL,
			body TEXT,
			occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// Indexes for common query patterns
		`CREATE INDEX IF NOT EXISTS idx_contacts_owner ON contacts(owner_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deals_stage ON deals(stage)`,
		`CREATE INDEX IF NOT EXISTS idx_deals_contact ON deals(contact_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activities_contact ON activities(contact_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activities_deal ON activities(deal_id)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
