package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"time"
)

type Db struct {
	conn *sql.DB
}

func InitDB(configs map[string]string) *Db {
	port, _ := strconv.Atoi(configs["SQL_PORT"])
	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configs["SQL_ADDRESS"], port, configs["SQL_USER"], configs["SQL_PASSWORD"], configs["SQL_DB_NAME"])

	db, err := sql.Open("postgres", info)
	if err != nil {
		log.Fatal(err)
	}

	// Tune the pool:
	db.SetConnMaxLifetime(5 * time.Minute)

	// Generate tables that aren't there
	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS lists (
        	id SERIAL PRIMARY KEY,
		    name TEXT UNIQUE NOT NULL
    	);
		ALTER TABLE lists 
		ADD COLUMN IF NOT EXISTS recipients TEXT[];
	`); err != nil {
		log.Fatal(err)
	}

	return &Db{db}
}
