// Package db provides functions for initializing and interacting with a database.
package db

import (
	"database/sql"
	"embed"
	"log"
	"os"
	"strings"
	_ "modernc.org/sqlite"
)

// initSQL is an embedded file system that contains the SQL script for initializing the database.
//
//go:embed init.sql
var initSQL embed.FS

// InitializeDatabase initializes the database by creating the database file if it doesn't exist,
// opening a connection to the database, and executing the SQL statements from the init.sql file.
// It returns an error if any error occurs during the initialization process.
func InitializeDatabase() error {
	dbFilePath := strings.Split(os.Getenv("DATABASE_URL"), ":")[1]

	// Check if the database file exists
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		// Create the database file if it doesn't existrc
		log.Print("Creating database file....")
		file, err := os.Create(dbFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		log.Printf("Database file(s) is already exists.")
		return nil
	}

	// Open a connection to the database
	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return err
	} else {
		defer db.Close()
	}

	// Read the SQL statements from the init.sql file
	initSqlFile, err := initSQL.ReadFile("init.sql")
	sqlStatements := strings.Split(string(initSqlFile), ";")
	for _, statement := range sqlStatements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		} else {
			statement += ";"
		}

		// Execute each SQL statement
		log.Print("Exec query: ", statement)
		_, err = db.Exec(statement)
		if err != nil {
			log.Fatalf("Error executing statement: %s\n%v", statement, err)
		}
	}

	return nil
}
