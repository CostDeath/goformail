package db

import (
	"github.com/lib/pq"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest/model"
)

func (db *Db) GetList(id int) (*model.List, *Error) {
	var list model.List
	if err := db.conn.QueryRow(`
		SELECT name, recipients FROM lists WHERE id = $1
	`, id,
	).Scan(&list.Name, pq.Array(&list.Recipients)); err != nil {
		return nil, getError(err)
	}

	return &list, nil
}

func (db *Db) CreateList(name string, recipients []string) (int, *Error) {
	var id int
	if err := db.conn.QueryRow(`
		INSERT INTO lists (name, recipients) VALUES ($1, $2) RETURNING id
	`, name, pq.Array(recipients),
	).Scan(&id); err != nil {
		return 0, getError(err)
	}

	return id, nil
}

func (db *Db) PatchList(id int, name string, recipients []string) *Error {
	if _, err := db.conn.Exec(`
		UPDATE lists
		SET
		    name = COALESCE(NULLIF($1, ''), name),
		    recipients = COALESCE(NULLIF($2, ARRAY[]::TEXT[]), recipients)
		WHERE id = $3;
	`, name, pq.Array(recipients), id,
	); err != nil {
		return getError(err)
	}

	return nil
}

func (db *Db) DeleteList(id int) *Error {
	if _, err := db.conn.Exec(`
		DELETE FROM lists WHERE id = $1;
	`, id,
	); err != nil {
		return getError(err)
	}

	return nil
}

func (db *Db) GetAllLists() (*[]model.List, *Error) {
	var lists []model.List
	rows, err := db.conn.Query(`
		SELECT name, recipients FROM lists
	`)

	for rows.Next() {
		var list model.List
		if err := rows.Scan(&list.Name, pq.Array(&list.Recipients)); err != nil {
			return nil, getError(err)
		}

		lists = append(lists, list)
	}

	if err != nil {
		return nil, getError(err)
	}

	return &lists, nil
}

func (db *Db) GetRecipientsFromListName(name string) ([]string, error) {
	var recipients []string
	if err := db.conn.QueryRow(`
		SELECT recipients FROM lists WHERE name = $1
	`, name,
	).Scan(&recipients); err != nil {
		return nil, err
	}

	return recipients, nil
}
