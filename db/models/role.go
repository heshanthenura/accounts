package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sliitmozilla/accounts/db"
	apiErrors "github.com/sliitmozilla/accounts/errors"
)

type RoleModel struct {
	Name string `json:"name"`
}

func (RoleModel) SelectAll() (roles []string, err error) {
	conn, err := db.ConnectDB()
	if err != nil {
		return
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(),
		"SELECT * FROM roles",
	)

	roles, err = pgx.CollectRows(rows, pgx.RowTo[string])
	return
}

func (r *RoleModel) Insert() (int, error) {
	conn, err := db.ConnectDB()
	if err != nil {
		return 0, nil
	}
	defer conn.Close(context.Background())

	t, err := conn.Exec(context.Background(),
		"INSERT INTO roles VALUES ($1)",
		r.Name,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, apiErrors.DuplicateError{Msg: "Role already exists"}
			}
		}
		return 0, err
	}
	return int(t.RowsAffected()), nil
}

func (r *RoleModel) Update(newRole RoleModel) (int, error) {
	conn, err := db.ConnectDB()
	if err != nil {
		return 0, err
	}
	defer conn.Close(context.Background())

	t, err := conn.Exec(context.Background(),
		"UPDATE roles SET name = $1 WHERE name = $2",
		newRole.Name, r.Name,
	)
	if err != nil {
		return 0, err
	}
	return int(t.RowsAffected()), err
}

func (r *RoleModel) Delete() (int, error) {
	conn, err := db.ConnectDB()
	if err != nil {
		return 0, err
	}
	defer conn.Close(context.Background())

	t, err := conn.Exec(context.Background(),
		"DELETE FROM roles WHERE name = $1",
		r.Name,
	)
	if err != nil {
		return 0, err
	}
	return int(t.RowsAffected()), err
}
