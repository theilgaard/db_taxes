package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	"theilgaard/db_taxes/internal/db"
)

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

func getTaxRecords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func getTaxRecordsByMunicipalityAndDate(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Both municipality and date specified, return single record for that municipality and date
		// Precedence on period_type is made by Daily, Weekly, Monthly, Yearly as 1, 2, 3, 4 respectively.
		municipality := c.Param("municipality")
		date := c.Query("date")
		if date == "" {
			// Only municipality specified, return all records for that municipality
			query := "SELECT * FROM tax_records WHERE municipality = ?;"
			rows, err := db.Query(query, municipality)
			if err != nil {
				log.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch database records"})
				return
			}
			parseRows(c, rows)
			return
		}
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
}

func postTaxRecord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the JSON body
		var record TaxRecord
		if err := c.BindJSON(&record); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert the record into the database
		if err := insertTaxRecord(db, record); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

func insertTaxRecord(db *sql.DB, record TaxRecord) error {
	const insert string = `
		INSERT INTO tax_records (municipality, period_type, date_start, date_end, tax_rate)
		VALUES (?, ?, ?, ?, ?);`
	_, err := db.Exec(insert, record.Municipality, record.PeriodType, record.DateStart, record.DateEnd, record.TaxRate)
	return err
}

func configureServer(db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/records", getTaxRecords(db))
	router.GET("/records/:municipality", getTaxRecordsByMunicipalityAndDate(db))

	router.POST("/records", postTaxRecord(db))

	return router
}
