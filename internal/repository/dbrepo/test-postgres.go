package dbrepo

import (
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"time"
)

func (p *testPostgresDBRepo) AllUsers() bool {
	return true
}

func (p *testPostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil

}

func (p *testPostgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

// SearchAvailabilityByDatesByRoomID return true if there is availability rooms exits, and false if no availability exits
func (p *testPostgresDBRepo) SearchAvailabilityByDatesByRoomID(roomId int, start, end time.Time) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any , for given data range
func (p *testPostgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

// GetRoomByID get room by id
func (p *testPostgresDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	return room, nil
}
