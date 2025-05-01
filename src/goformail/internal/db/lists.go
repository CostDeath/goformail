package db

import (
	"fmt"
	"github.com/lib/pq"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
)

func (db *Db) GetList(id int) (*model.ListResponse, *Error) {
	list := model.ListResponse{Id: id}
	if err := db.conn.QueryRow(`
		SELECT name, recipients, mods, approved_senders, locked FROM lists WHERE id = $1
	`, id,
	).Scan(&list.Name, pq.Array(&list.Recipients), pq.Array(&list.Mods), pq.Array(&list.ApprovedSenders),
		&list.Locked); err != nil {
		return nil, getError(err)
	}

	return &list, nil
}

func (db *Db) CreateList(list *model.ListRequest) (int, *Error) {
	var id int
	if err := db.conn.QueryRow(`
		INSERT INTO lists (name, recipients, mods, approved_senders, locked) VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, list.Name, pq.Array(list.Recipients), pq.Array(list.Mods), pq.Array(list.ApprovedSenders), list.Locked,
	).Scan(&id); err != nil {
		return 0, getError(err)
	}

	return id, nil
}

func (db *Db) PatchList(id int, list *model.ListRequest, override *model.ListOverrides) *Error {
	rcpt := `recipients = COALESCE(NULLIF($2, ARRAY[]::TEXT[]), recipients)`
	mods := `mods = COALESCE(NULLIF($3, ARRAY[]::INT[]), mods)`
	senders := `approved_senders = COALESCE(NULLIF($4, ARRAY[]::TEXT[]), approved_senders)`
	locked := "locked = locked AND $5 IS NOT DISTINCT FROM $5" // forced to use $5 for it to run
	if override != nil && override.Recipients {
		rcpt = `recipients = $2`
	}
	if override != nil && override.Mods {
		mods = `mods = $3`
	}
	if override != nil && override.ApprovedSenders {
		senders = `approved_senders = $4`
	}
	if list.Locked || override != nil && override.Locked {
		locked = `locked = $5`
	}

	res, err := db.conn.Exec(fmt.Sprintf(`
		UPDATE lists
		SET
		    name = COALESCE(NULLIF($1, ''), name), %s, %s, %s, %s
		WHERE id = $6;
	`, rcpt, mods, senders, locked),
		list.Name, pq.Array(list.Recipients), pq.Array(list.Mods), pq.Array(list.ApprovedSenders), list.Locked, id,
	)
	if err != nil {
		return getError(err)
	} else if count, err := res.RowsAffected(); count == 0 || err != nil {
		return &Error{Err: err, Code: ErrNoRows}
	}

	return nil
}

func (db *Db) DeleteList(id int) *Error {
	res, err := db.conn.Exec(`
		DELETE FROM lists WHERE id = $1;
	`, id,
	)
	if err != nil {
		return getError(err)
	} else if count, err := res.RowsAffected(); count == 0 || err != nil {
		return &Error{Err: err, Code: ErrNoRows}
	}

	return nil
}

func (db *Db) GetAllLists() (*[]*model.ListResponse, *Error) {
	lists := []*model.ListResponse{}
	rows, err := db.conn.Query(`
		SELECT id, name, recipients, mods, approved_senders, locked FROM lists
	`)
	if err != nil {
		return nil, getError(err)
	}

	for rows.Next() {
		list := model.ListResponse{}
		if err := rows.Scan(&list.Id, &list.Name, pq.Array(&list.Recipients), pq.Array(&list.Mods),
			pq.Array(&list.ApprovedSenders), &list.Locked); err != nil {
			return nil, getError(err)
		}

		lists = append(lists, &list)
	}

	return &lists, nil
}

func (db *Db) GetApprovalFromListName(sender string, name string) (int, bool, *Error) {
	id := 0
	approved := false
	if err := db.conn.QueryRow(`
		SELECT id, $1 = ANY(approved_senders)
		FROM lists WHERE name = $2;
	`, sender, name,
	).Scan(&id, &approved); err != nil {
		return 0, false, getError(err)
	}

	return id, approved, nil
}
