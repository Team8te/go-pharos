package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	deviceTable = "device"
)

type Repository struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func NewRepository(db *sqlx.DB) *Repository {
	schema := `
	CREATE TABLE device (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		uuid text,
		ip text,
		create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(schema)
	if err != nil {
		panic(err)
	}

	return &Repository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
