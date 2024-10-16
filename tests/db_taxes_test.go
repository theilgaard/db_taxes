package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"theilgaard/db_taxes/internal/db"
	"theilgaard/db_taxes/cmd/db_taxes" // Import the configureServer function
)

func setupTestServer() *httptest.Server {
	database, _ := db.InitializeDatabase()
	db.PopulateDatabase(database)
	return httptest.NewServer(configureServer(database))
}

func TestAPIEndpoints(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	test_url := "/records/Copenhagen?date=2024-01-01"
	expected_tax_rate := 0.1
	t.Run("Success GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []db.TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, expected_tax_rate, tax_records[0].TaxRate)
	})

	test_url = "/records/Copenhagen?date=2024-03-16"
	expected_tax_rate = 0.2
	t.Run("Success GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []db.TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, expected_tax_rate, tax_records[0].TaxRate)
	})

	test_url = "/records/Copenhagen?date=2024-05-02"
	expected_tax_rate = 0.4
	t.Run("Success GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []db.TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, expected_tax_rate, tax_records[0].TaxRate)
	})

	test_url = "/records/Copenhagen?date=2024-07-10"
	expected_tax_rate = 0.2
	t.Run("Success GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []db.TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, expected_tax_rate, tax_records[0].TaxRate)
	})

	test_url = "/records/Copenhagen?date=invalid_date_format"
	expected_tax_rate = 0.2
	t.Run("Invalid Date Format GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []db.TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Contains(t, string(response_bytes), "Invalid date format")
	})

	test_url = "/records/"
	kolding_start_date, _ := time.Parse(time.DateOnly, "2024-01-01")
	kolding_end_date, _ := time.Parse(time.DateOnly, "2024-12-31")
	kolding_tax_record := db.TaxRecord{Municipality: "Kolding", PeriodType: 4, DateStart: kolding_start_date, DateEnd: kolding_end_date, TaxRate: 0.2}
	kolding_tax_record_json, _ := json.Marshal(kolding_tax_record)
	t.Run("Success POST "+test_url, func(t *testing.T) {
		resp, err := http.Post(ts.URL+test_url, "application/json", bytes.NewBuffer(kolding_tax_record_json))
		if err != nil {
			t.Fatalf("Failed to send POST request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}

func TestInsertAndReadRecord(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	test_url := "/records/"
	new_record := db.TaxRecord{
		Municipality: "TestCity",
		PeriodType:   3,
		DateStart:    time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		DateEnd:      time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
		TaxRate:      0.5,
	}
	new_record_json, _ := json.Marshal(new_record)

	t.Run("Insert Record", func(t *testing.T) {
		resp, err := http.Post(ts.URL+test_url, "application/json", bytes.NewBuffer(new_record_json))
		if err != nil {
			t.Fatalf("Failed to send POST request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Read Inserted Record", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/records/TestCity?date=2024-06-15")
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []db.TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, 1, len(tax_records))
		assert.Equal(t, new_record.Municipality, tax_records[0].Municipality)
		assert.Equal(t, new_record.PeriodType, tax_records[0].PeriodType)
		assert.Equal(t, new_record.DateStart, tax_records[0].DateStart)
		assert.Equal(t, new_record.DateEnd, tax_records[0].DateEnd)
		assert.Equal(t, new_record.TaxRate, tax_records[0].TaxRate)
	})
}
