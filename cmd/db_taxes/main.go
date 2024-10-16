package main

import (
	"log"
	"theilgaard/db_taxes/internal/db"
)

func main() {
	// Initialize the database
	database, err := db.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Populate the database with initial records
	if err := db.PopulateDatabase(database); err != nil {
		log.Fatalf("Failed to populate database: %v", err)
	}

	// Run the server
	router := configureServer(database)
	router.Run(":8080")
}
