package database

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/MyCloudNest/Nest-Server/utils"
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
	schemaPath := filepath.Join(dbDir, "schema.sql")

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

	schemaURL := "https://raw.githubusercontent.com/MyCloudNest/Nest-Server/refs/heads/main/database/migrations/schema.sql"
	if _, err = os.Stat(schemaPath); os.IsNotExist(err) {
		log.Println("Schema file not found. Downloading...")
		if err = utils.DownloadFile(schemaURL, schemaPath); err != nil {
			log.Fatalf("Failed to download schema file: %v", err)
		}
		log.Println("Schema file downloaded successfully:", schemaPath)
	}

	if err = executeSQLFile(schemaPath); err != nil {
		log.Printf("Warning: Failed to execute some schema statements: %v", err)
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

func executeSQLFile(filePath string) error {
	sqlFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer sqlFile.Close()

	sqlContent, err := io.ReadAll(sqlFile)
	if err != nil {
		return err
	}

	statements := strings.Split(string(sqlContent), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err := DB.Exec(stmt)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("Skipping already existing object: %s", stmt)
				continue
			}
			return err
		}
	}

	return nil
}
