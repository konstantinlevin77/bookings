package repository

import "github.com/konstantinlevin77/bookings/internal/models"

type DatabaseRepo interface {
	InsertReservation(r models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
}
