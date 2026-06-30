package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
)

func (p *testPostgresDBRepo) AllUsers() bool {
	return true
}

func (p *testPostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID > 2 {
		return 1, errors.New("some error")
	}
	return 1, nil

}

func (p *testPostgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 2 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID return true if there is availability rooms exits, and false if no availability exits
func (p *testPostgresDBRepo) SearchAvailabilityByDatesByRoomID(roomId int, start, end time.Time) (bool, error) {
	// set up a test time
	layout := "2006-01-02"
	str := "2049-12-31"
	testDateToFail, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}
	if start == testDateToFail {
		return false, errors.New("error querying database")
	}
	if start.After(testDateToFail) {
		return false, nil
	}
	return true, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any , for given data range
func (p *testPostgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	endToFail, _ := time.Parse("2006-01-02", "2050-01-03")
	if endToFail == end {
		return []models.Room{}, errors.New("some error")
	}
	lenFail, _ := time.Parse("2006-01-02", "2050-01-04")
	if lenFail == end {
		return []models.Room{}, nil
	}
	var rooms []models.Room
	rooms = append(rooms, models.Room{})
	return rooms, nil

}

// GetUserByID returns a user by ID for testing
func (p *testPostgresDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	if id > 10 {
		return user, errors.New("user not found")
	}
	user.ID = id
	user.FirstName = "Test"
	user.LastName = "User"
	user.Email = "test@example.com"
	user.AccessLevel = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user, nil
}

// GetRoomByID returns a room by ID for testing
func (p *testPostgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 100 {
		return room, errors.New("room not found")
	}
	room.ID = id
	room.RoomName = "Test Room"
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()
	return room, nil
}

// UpdaterUser updates a user for testing
func (p *testPostgresDBRepo) UpdaterUser(u models.User) error {
	if u.ID == 0 {
		return errors.New("user ID cannot be zero")
	}
	return nil
}

// Authenticate authenticates a user for testing
func (p *testPostgresDBRepo) Authenticate(email, password string) (int, string, error) {
	if email == "fail@example.com" {
		return 0, "", errors.New("authentication failed")
	}
	if email == "admin@example.com" {
		return 1, "admin", nil
	}
	return 2, "user", nil
}

func (p *testPostgresDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (p *testPostgresDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationByID returns a reservation by ID for testing
func (p *testPostgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation
	if id == 123 {
		return reservation, errors.New("reservation not found")
	}
	return reservation, nil
}

// UpdateReservation updates a reservation for testing
func (p *testPostgresDBRepo) UpdateReservation(reservation models.Reservation) error {
	if reservation.ID == 0 {
		return errors.New("reservation ID cannot be zero")
	}
	return nil
}

// DeleteReservation deletes a reservation for testing
func (p *testPostgresDBRepo) DeleteReservation(id int) error {
	if id == 0 {
		return errors.New("reservation ID cannot be zero")
	}
	return nil
}

// UpdateProcessedForReservation updates the processed status of a reservation for testing
func (p *testPostgresDBRepo) UpdateProcessedForReservation(id int, processed int) error {
	if id == 0 {
		return errors.New("reservation ID cannot be zero")
	}
	return nil
}
