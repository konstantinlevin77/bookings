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

func (pr *postgresRepo) InsertReservation(r models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	var query string
	query = `insert into reservations (first_name, last_name, email, phone, start_date,
             end_date, room_id, created_at, updated_at) values   ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := pr.DB.QueryRowContext(ctx, query,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	return newID, err
}

func (pr *postgresRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id,
                               created_at, updated_at, restriction_id)
                               VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := pr.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	return err
}

// SearchAvailabilityByRoomID checks whether a date range is available or not.
func (pr *postgresRepo) SearchAvailabilityByRoomID(startDate, endDate time.Time, roomID int) (bool, error) {

	var numRows int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			SELECT COUNT(id) FROM room_restrictions 
            WHERE room_id = $1 AND 
            $2 < end_date and $3 > start_date`

	row := pr.DB.QueryRowContext(ctx, query, roomID, startDate, endDate)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms checks if any room is available and returns them if they exist.
func (pr *postgresRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
				SELECT r.id r.room_name 
				FROM rooms r WHERE r.id 
                NOT IN (SELECT room_id FROM room_restrictions rr WHERE $1 < rr.end_date AND $2 > rr.start_date)`

	rows, err := pr.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresRepo{
		App: a,
		DB:  conn,
	}
}
