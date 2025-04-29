package db

import (
	"github.com/lib/pq"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
)

func (db *Db) GetUser(id int) (*model.UserResponse, *Error) {
	user := model.UserResponse{Id: id}
	if err := db.conn.QueryRow(`
		SELECT email, permissions FROM users WHERE id = $1
	`, id,
	).Scan(&user.Email, pq.Array(&user.Permissions)); err != nil {
		return nil, getError(err)
	}

	return &user, nil
}

func (db *Db) CreateUser(user *model.UserRequest, hash string) (int, *Error) {
	var id int
	if err := db.conn.QueryRow(`
		INSERT INTO users (email, hash, permissions) VALUES ($1, $2, $3) RETURNING id
	`, user.Email, hash, pq.Array(user.Permissions),
	).Scan(&id); err != nil {
		return 0, getError(err)
	}

	return id, nil
}

func (db *Db) UpdateUser(id int, user *model.UserRequest) *Error {
	res, err := db.conn.Exec(`
		UPDATE users
		SET
		    email = COALESCE(NULLIF($1, ''), email),
		    permissions = COALESCE(NULLIF($2, ARRAY[]::TEXT[]), permissions)
		WHERE id = $3;
	`, user.Email, pq.Array(user.Permissions), id,
	)
	if err != nil {
		return getError(err)
	} else if count, err := res.RowsAffected(); count == 0 || err != nil {
		return &Error{Err: err, Code: ErrNoRows}
	}

	return nil
}

func (db *Db) DeleteUser(id int) *Error {
	res, err := db.conn.Exec(`
		DELETE FROM users WHERE id = $1;
	`, id,
	)
	if err != nil {
		return getError(err)
	} else if count, err := res.RowsAffected(); count == 0 || err != nil {
		return &Error{Err: err, Code: ErrNoRows}
	}

	return nil
}

func (db *Db) GetAllUsers() (*[]*model.UserResponse, *Error) {
	var users []*model.UserResponse
	rows, err := db.conn.Query(`
		SELECT id, email, permissions FROM users
	`)
	if err != nil {
		return nil, getError(err)
	}

	for rows.Next() {
		user := model.UserResponse{}
		if err := rows.Scan(&user.Id, &user.Email, pq.Array(&user.Permissions)); err != nil {
			return nil, getError(err)
		}

		users = append(users, &user)
	}

	return &users, nil
}

func (db *Db) GetUserPassword(email string) (int, string, *Error) {
	var id int
	var hash string
	if err := db.conn.QueryRow(`
		SELECT id, hash FROM users WHERE email = $1
	`, email,
	).Scan(&id, &hash); err != nil {
		return 0, "", getError(err)
	}

	return id, hash, nil
}

func (db *Db) UserExists(id int) (bool, *Error) {
	exists := false
	if err := db.conn.QueryRow(`
		SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)
	`, id,
	).Scan(&exists); err != nil {
		return exists, getError(err)
	}

	return exists, nil
}

func (db *Db) GetUserPerms(id int) ([]string, *Error) {
	var perms []string
	if err := db.conn.QueryRow(`
		SELECT permissions FROM users WHERE id = $1
	`, id,
	).Scan(pq.Array(&perms)); err != nil {
		return nil, getError(err)
	}

	return perms, nil
}
