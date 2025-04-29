package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"log"
	"strconv"
	"time"
)

type IDb interface {
	GetList(id int) (*model.List, *Error)
	CreateList(list *model.List) (int, *Error)
	PatchList(id int, list *model.List) *Error
	DeleteList(id int) *Error
	GetAllLists() (*[]*model.ListWithId, *Error)
	GetRecipientsFromListName(name string) ([]string, error)
	GetUser(id int) (*model.UserResponse, *Error)
	CreateUser(user *model.UserRequest, hash string, salt string) (int, *Error)
	UpdateUser(id int, user *model.UserRequest) *Error
	DeleteUser(id int) *Error
	GetAllUsers() (*[]*model.UserResponse, *Error)
}

type Db struct {
	IDb
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
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS lists (
        	id SERIAL PRIMARY KEY
    	);
		ALTER TABLE lists
		    ADD COLUMN IF NOT EXISTS name TEXT UNIQUE NOT NULL,
		    ADD COLUMN IF NOT EXISTS recipients TEXT[];

		CREATE TABLE IF NOT EXISTS users (
        	id SERIAL PRIMARY KEY
    	);
		ALTER TABLE users
		    ADD COLUMN IF NOT EXISTS email TEXT UNIQUE NOT NULL,
		    ADD COLUMN IF NOT EXISTS hash TEXT NOT NULL,
		    ADD COLUMN IF NOT EXISTS salt TEXT NOT NULL,
		    ADD COLUMN IF NOT EXISTS token TEXT,
		    ADD COLUMN IF NOT EXISTS permissions TEXT[];
	`); err != nil {
		log.Fatal(err)
	}

	return &Db{conn: db}
}
