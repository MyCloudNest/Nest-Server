package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	dbDir := filepath.Join(homeDir, ".cloudnest")
	dbPath := filepath.Join(dbDir, "db.sqlite")

	if err = os.MkdirAll(dbDir, 0o700); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if DB == nil {
		log.Fatal(errors.New("database connection is nil"))
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	_, err = DB.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA temp_store = MEMORY;
		PRAGMA cache_size = -64000;
		PRAGMA mmap_size = 30000000000;
		PRAGMA optimize;
	`)
	if err != nil {
		log.Fatalf("Failed to set PRAGMA statements: %v", err)
	}

	log.Println("Database initialized successfully")
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}
}
