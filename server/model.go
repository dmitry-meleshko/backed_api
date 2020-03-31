package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type patientItem struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type patients struct {
	items []patientItem
}

// unused for now
type physicianItem struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
}

// unused for now
type visitItem struct {
	PatientID   int    `json:"patient_id"`
	PhysicianID int    `json:"physician_id"`
	VisitedAt   string `json:"visited_at"`
	Location    string `json:"location"`
	Reason      string `json:"reason"`
}

func (p *patientItem) getPatient(db *sql.DB) error {
	err := db.QueryRow("SELECT first_name, last_name, address, phone, email, created_at "+
		"FROM patient WHERE id=$1",
		p.ID).Scan(&p.FirstName, &p.LastName, &p.Address, &p.Phone, &p.Email, &p.CreatedAt)

	return err
}

func (p *patientItem) addPatient(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO patient(first_name, last_name, address, phone, email) "+
		"VALUES($1, $2, $3, $4, $5) RETURNING id",
		p.FirstName, p.LastName, p.Address, p.Phone, p.Email).Scan(&p.ID)

	return err
}

func (ps *patients) getPatients(db *sql.DB) ([]patientItem, error) {
	// TODO: convert to channel for feeding data
	rows, err := db.Query("SELECT id, first_name, last_name, address, phone, email, created_at FROM patient")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p patientItem
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Address, &p.Phone, &p.Email, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		// TODO: feed to a channel here
		ps.items = append(ps.items, p)
	}

	return ps.items, nil
}

func (p *patientItem) deletePatient(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM patient WHERE id = $1", p.ID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()

	return err
}

// TODO: implement physicianItem and visitItem methods
