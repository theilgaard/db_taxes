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
)

func setupTestServer() *httptest.Server {
	db, _ := initializeDatabase()
	populateDatabase(db)
	return httptest.NewServer(configureServer(db))
}

func TestAPIEndpoints(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	test_url := "/records/Copenhagen/2024-01-01"
	expected_tax_rate := 0.1
	t.Run("Success GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, expected_tax_rate, tax_records[0].TaxRate)
	})

	test_url = "/records/Copenhagen/2024-03-16"
	expected_tax_rate = 0.2
	t.Run("Success GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, expected_tax_rate, tax_records[0].TaxRate)
	})

	test_url = "/records/Copenhagen/2024-05-02"
	expected_tax_rate = 0.4
	t.Run("Success GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, expected_tax_rate, tax_records[0].TaxRate)
	})

	test_url = "/records/Copenhagen/2024-07-10"
	expected_tax_rate = 0.2
	t.Run("Success GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Equal(t, expected_tax_rate, tax_records[0].TaxRate)
	})

	test_url = "/records/Copenhagen/invalid_date_format"
	expected_tax_rate = 0.2
	t.Run("Invalid Date Format GET "+test_url, func(t *testing.T) {
		resp, err := http.Get(ts.URL + test_url)
		if err != nil {
			t.Fatalf("Failed to send GET request: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		response_bytes, _ := io.ReadAll(resp.Body)
		tax_records := []TaxRecord{}
		json.Unmarshal(response_bytes, &tax_records)

		assert.Contains(t, string(response_bytes), "Invalid date format")
	})

	test_url = "/records/"
	kolding_start_date, _ := time.Parse(time.DateOnly, "2024-01-01")
	kolding_end_date, _ := time.Parse(time.DateOnly, "2024-12-31")
	kolding_tax_record := TaxRecord{Municipality: "Kolding", PeriodType: 4, DateStart: kolding_start_date, DateEnd: kolding_end_date, TaxRate: 0.2}
	kolding_tax_record_json, _ := json.Marshal(kolding_tax_record)
	t.Run("Success POST "+test_url, func(t *testing.T) {
		resp, err := http.Post(ts.URL+test_url, "application/json", bytes.NewBuffer(kolding_tax_record_json))
		if err != nil {
			t.Fatalf("Failed to send POST request: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}
