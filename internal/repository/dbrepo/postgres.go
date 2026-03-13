package dbrepo

import (
	"context"
	"time"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
)

func (p *postgresDBRepo) AllUsers() bool {
	return true
}

func (p *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	var newId int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
		end_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := p.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, err

}

func (p *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
		restriction_id, created_at, updated_at)
		values
		($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID return true if there is availability rooms exits, and false if no availability exits
func (p *postgresDBRepo) SearchAvailabilityByDatesByRoomID(roomId int, start, end time.Time) (bool, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	var numRows int
	stmt := `
		SELECT
			COUNT("id")
		FROM
			room_restrictions
		WHERE
			room_id = $1
			AND end_date > $2 AND start_date < $3;
	`
	err := p.DB.QueryRowContext(ctx, stmt,
		roomId,
		start,
		end,
	).Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any , for given data range
func (p *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	var rooms []models.Room
	stmt := `
		select r.id ,r.room_name from rooms r 
		where r.id not in (
			select rr.id from room_restrictions rr 
				where $1 < rr.end_date 
				and $2 >rr.start_date
			)
	`
	rows, err := p.DB.QueryContext(ctx, stmt,
		start,
		end,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}

// GetRoomByID get room by id
func (p *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*3)
	defer cancle()

	var room models.Room
	stmt := `
		select id, room_name, created_at, updated_at from rooms where id = $1
	`

	row := p.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}
