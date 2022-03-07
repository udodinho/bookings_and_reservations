package dbrepo

import (
	"database/sql"
	"github.com/udodinho/bookings/internal/config"
	"github.com/udodinho/bookings/repository"
)

type PostgresDbRepo struct {
	DB  *sql.DB
	App *config.AppConfig
}

// NewPostgresDbRepo creates a new instance of PostgresDbRepo
func NewPostgresDbRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDbRepo{
		App: a,
		DB:  conn,
	}
}
