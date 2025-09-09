package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	db "server/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func NewDatabase(schemaPath string) (*db.Queries, *sql.DB) {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("[DATABASE] No .env file found, using system environment variables")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("[DATABASE] Missing DB_URL environment variable")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("[DATABASE] Failed to open DB: %v", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatalf("[DATABASE] Failed to ping DB: %v", err)
	}
	log.Println("[DATABASE] Connected to PostgreSQL")

	if schemaPath != "" {
		if err := applySchema(conn, schemaPath); err != nil {
			log.Fatalf("[DATABASE] Failed to apply schema: %v", err)
		}
		log.Println("[DATABASE] Schema applied successfully")
	}

	return db.New(conn), conn
}

func applySchema(conn *sql.DB, path string) error {
	schemaSQL, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	_, err = conn.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}
