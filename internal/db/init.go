package db

import (
	"database/sql"
	"time"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

const db_driver string = "sqlite3"
const db_location string = "tax_records.db"

type TaxRecord struct {
	Municipality string    `json:"municipality"`
	PeriodType   int       `json:"period_type"`
	DateStart    time.Time `json:"date_start"`
	DateEnd      time.Time `json:"date_end"`
	TaxRate      float64   `json:"tax_rate"`
}

func getInitialTaxRecords() []TaxRecord {
	yearly_start_time, _ := time.Parse(time.DateOnly, "2024-01-01")
	yearly_end_time, _ := time.Parse(time.DateOnly, "2024-12-31")
	monthly_start_time, _ := time.Parse(time.DateOnly, "2024-05-01")
	monthly_end_time, _ := time.Parse(time.DateOnly, "2024-05-31")
	daily_1_start_time, _ := time.Parse(time.DateOnly, "2024-01-01")
	daily_1_end_time, _ := time.Parse(time.DateOnly, "2024-01-01")
	daily_2_start_time, _ := time.Parse(time.DateOnly, "2024-12-25")
	daily_2_end_time, _ := time.Parse(time.DateOnly, "2024-12-25")

	var tax_records = []TaxRecord{
		{Municipality: "Copenhagen", PeriodType: 4, DateStart: yearly_start_time, DateEnd: yearly_end_time, TaxRate: 0.2},
		{Municipality: "Copenhagen", PeriodType: 3, DateStart: monthly_start_time, DateEnd: monthly_end_time, TaxRate: 0.4},
		{Municipality: "Copenhagen", PeriodType: 1, DateStart: daily_1_start_time, DateEnd: daily_1_end_time, TaxRate: 0.1},
		{Municipality: "Copenhagen", PeriodType: 1, DateStart: daily_2_start_time, DateEnd: daily_2_end_time, TaxRate: 0.1},
		{Municipality: "Aarhus", PeriodType: 4, DateStart: yearly_start_time, DateEnd: yearly_end_time, TaxRate: 0.5},
	}
	return tax_records
}

func initializeDatabase() (*sql.DB, error) {
	// Initialize the database with some initial records
	db, err := sql.Open(db_driver, db_location)
	if (err != nil) {
		return nil, err
	}

	const drop string = "DROP TABLE IF EXISTS tax_records;"
	if _, err := db.Exec(drop); err != nil {
		return nil, err
	}

	const create string = `
		CREATE TABLE tax_records (
			id INTEGER NOT NULL PRIMARY KEY,
			municipality TEXT NOT NULL,
			period_type INT NOT NULL,
			date_start DATETIME NOT NULL,
			date_end DATETIME NOT NULL,
			tax_rate REAL NOT NULL
	);`
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return db, nil
}

func insertTaxRecord(db *sql.DB, record TaxRecord) error {
	const insert string = `
		INSERT INTO tax_records (municipality, period_type, date_start, date_end, tax_rate)
		VALUES (?, ?, ?, ?, ?);`
	_, err := db.Exec(insert, record.Municipality, record.PeriodType, record.DateStart, record.DateEnd, record.TaxRate)
	return err
}

func populateDatabase(db *sql.DB) error {
	// Insert the initial records into the database
	for _, record := range getInitialTaxRecords() {
		if err := insertTaxRecord(db, record); err != nil {
			return err
		}
	}
	return nil
}
