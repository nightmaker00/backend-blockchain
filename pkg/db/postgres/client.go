package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode)
	log.Printf("Connecting to database with config: host=%s port=%s user=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	// Проверяем схему
	var schemaName string
	err = db.Get(&schemaName, "SELECT current_schema()")
	if err != nil {
		log.Printf("Error getting current schema: %v", err)
	} else {
		log.Printf("Current schema: %s", schemaName)
	}

	// Проверяем таблицы
	var tables []string
	err = db.Select(&tables, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		log.Printf("Error getting tables: %v", err)
	} else {
		log.Printf("Available tables: %v", tables)
	}

	return db, nil
}
