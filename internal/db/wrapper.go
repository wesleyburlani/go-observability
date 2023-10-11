package db

import (
	"database/sql"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Database struct {
	DB      *sql.DB
	Queries *Queries
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := otelsql.Open("postgres", dsn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		return nil, err
	}

	queries := New(db)
	return &Database{
		DB:      db,
		Queries: queries,
	}, nil
}
