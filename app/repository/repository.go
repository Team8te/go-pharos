package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

const (
	deviceTable = "device"
)

type Repository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
