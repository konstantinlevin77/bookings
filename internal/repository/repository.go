package repository

import (
	"github.com/konstantinlevin77/bookings/internal/models"
	"time"
)

type DatabaseRepo interface {
	InsertReservation(r models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByRoomID(startDate, endDate time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error)
}
