package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"log"
	"strconv"
	"time"
)

type IDb interface {
	GetList(id int) (*model.ListResponse, *Error)
	CreateList(list *model.ListRequest) (int, *Error)
	PatchList(id int, list *model.ListRequest, override *model.ListOverrides) *Error
	DeleteList(id int) *Error
	GetAllLists() (*[]*model.ListResponse, *Error)
	GetRecipientsFromListName(name string) ([]string, error)
	GetUser(id int) (*model.UserResponse, *Error)
	CreateUser(user *model.UserRequest, hash string) (int, *Error)
	UpdateUser(id int, user *model.UserRequest, overridePerms bool) *Error
	DeleteUser(id int) *Error
	GetAllUsers() (*[]*model.UserResponse, *Error)
	GetUserPassword(email string) (int, string, *Error)
	UserExists(id int) (bool, *Error)
	UsersExist(ids []int64) ([]int64, *Error)
	GetUserPerms(id int) ([]string, *Error)
	GetUserPermsAndModStatus(id int, listId int) ([]string, bool, *Error)
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
		CREATE TABLE IF NOT EXISTS props (
        	id SERIAL PRIMARY KEY
    	);
		ALTER TABLE props
		    ADD COLUMN IF NOT EXISTS jwt_secret BYTEA;
		INSERT INTO props DEFAULT VALUES;

		CREATE TABLE IF NOT EXISTS lists (
        	id SERIAL PRIMARY KEY
    	);
		ALTER TABLE lists
		    ADD COLUMN IF NOT EXISTS name TEXT UNIQUE NOT NULL,
		    ADD COLUMN IF NOT EXISTS recipients TEXT[],
		    ADD COLUMN IF NOT EXISTS mods INT[],
		    ADD COLUMN IF NOT EXISTS approved_senders TEXT[],
		    ADD COLUMN IF NOT EXISTS locked BOOL DEFAULT false;

		CREATE TABLE IF NOT EXISTS users (
        	id SERIAL PRIMARY KEY
    	);
		ALTER TABLE users
		    ADD COLUMN IF NOT EXISTS email TEXT UNIQUE NOT NULL,
		    ADD COLUMN IF NOT EXISTS hash TEXT NOT NULL,
		    ADD COLUMN IF NOT EXISTS permissions TEXT[];
	`); err != nil {
		log.Fatal(err)
	}

	return &Db{conn: db}
}

func (db *Db) GetJwtSecret() *[]byte {
	secret := make([]byte, 32)
	err := db.conn.QueryRow("SELECT jwt_secret FROM props LIMIT 1").Scan(&secret)
	if err != nil {
		log.Fatal(err)
	}

	if secret == nil {
		// Generate new secret
		secret = make([]byte, 32)
		_, err = rand.Read(secret)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.conn.Exec("UPDATE props SET jwt_secret = $1", secret)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &secret
}
