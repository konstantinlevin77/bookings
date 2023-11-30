package dbrepo

import (
	"context"
	"database/sql"
	"github.com/konstantinlevin77/bookings/internal/config"
	"github.com/konstantinlevin77/bookings/internal/models"
	"github.com/konstantinlevin77/bookings/internal/repository"
	"time"
)

type postgresRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func (pr *postgresRepo) InsertReservation(r models.Reservation) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var query string
	query = `insert into reservations (first_name, last_name, email, phone, start_date,
             end_date, room_id, created_at, updated_at) values   ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err := pr.DB.ExecContext(ctx, query,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		time.Now(),
		time.Now(),
	)

	return err
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresRepo{
		App: a,
		DB:  conn,
	}
}
