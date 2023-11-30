package dbrepo

import (
	"database/sql"
	"github.com/konstantinlevin77/bookings/internal/config"
	"github.com/konstantinlevin77/bookings/internal/repository"
)

type postgresRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresRepo{
		App: a,
		DB:  conn,
	}
}
