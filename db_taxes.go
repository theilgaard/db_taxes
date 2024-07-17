package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func parseRows(c *gin.Context, rows *sql.Rows) {
	tax_records := []TaxRecord{}
	for rows.Next() {
		unused_id := 0
		var record TaxRecord
		if err := rows.Scan(&unused_id, &record.Municipality, &record.PeriodType, &record.DateStart, &record.DateEnd, &record.TaxRate); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch database records"})
			return
		}
		tax_records = append(tax_records, record)
	}
	defer rows.Close()
	c.IndentedJSON(http.StatusOK, tax_records)
}

func getTaxRecords(c *gin.Context) {
	db, err := sql.Open(db_driver, db_location)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// No municipality or date specified, return all records
	query := "SELECT * FROM tax_records;"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch database records"})
		return
	}
	parseRows(c, rows)

}

func getTaxRecordsByMunicipality(c *gin.Context) {
	// Only municipality specified, return all records for that municipality
	db, err := sql.Open(db_driver, db_location)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	municipality := c.Param("municipality")
	query := "SELECT * FROM tax_records WHERE municipality = ?;"
	rows, err := db.Query(query, municipality)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch database records"})
		return
	}
	parseRows(c, rows)
}

func getTaxRecordsByMunicipalityAndDate(c *gin.Context) {
	// Both municipality and date specified, return single record for that municipality and date
	// Precedence on period_type is made by Daily, Weekly, Monthly, Yearly as 1, 2, 3, 4 respectively.
	db, err := sql.Open(db_driver, db_location)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	municipality := c.Param("municipality")
	date := c.Param("date")
	date_time, err := time.Parse(time.DateOnly, date)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	query := "SELECT * FROM tax_records WHERE municipality = ? AND (date_start <= ? AND date_end >= ?) ORDER BY period_type LIMIT 1;"
	rows, err := db.Query(query, municipality, date_time, date_time)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch database records"})
		return
	}
	parseRows(c, rows)
}

func postTaxRecord(c *gin.Context) {
	// Parse the JSON body
	var record TaxRecord
	if err := c.BindJSON(&record); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the record into the database
	db, err := sql.Open(db_driver, db_location)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := insertTaxRecord(db, record); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func insertTaxRecord(db *sql.DB, record TaxRecord) error {
	const insert string = `
		INSERT INTO tax_records (municipality, period_type, date_start, date_end, tax_rate)
		VALUES (?, ?, ?, ?, ?);`
	_, err := db.Exec(insert, record.Municipality, record.PeriodType, record.DateStart, record.DateEnd, record.TaxRate)
	return err
}

func initializeDatabase() (*sql.DB, error) {
	// Initialize the database with some initial records
	db, err := sql.Open(db_driver, db_location)
	if err != nil {
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

func populateDatabase(db *sql.DB) error {
	// Insert the initial records into the database
	for _, record := range getInitialTaxRecords() {
		if err := insertTaxRecord(db, record); err != nil {
			return err
		}
	}
	return nil
}

func configureServer() *gin.Engine {
	router := gin.Default()
	router.GET("/records", getTaxRecords)
	router.GET("/records/:municipality", getTaxRecordsByMunicipality)
	router.GET("/records/:municipality/:date", getTaxRecordsByMunicipalityAndDate)
	router.POST("/records", postTaxRecord)

	return router
}

func main() {
	// Initialize the database
	db, err := initializeDatabase()
	if err != nil {
		panic(err)
	}

	// Populate the database with initial records
	if err := populateDatabase(db); err != nil {
		panic(err)
	}

	// Run the server
	router := configureServer()
	router.Run(":8080")
}
