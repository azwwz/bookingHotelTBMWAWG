package dbrepo

import (
	"errors"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"log"
	"time"
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

// GetRoomByID get room by id
func (p *testPostgresDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room
	if id > 2 {
		return room, errors.New("id can't more than 2")
	}

	return room, nil
}
