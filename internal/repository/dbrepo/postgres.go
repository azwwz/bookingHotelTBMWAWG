package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	var newId int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
		end_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
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

func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
		restriction_id, created_at, updated_at)
		values
		($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
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
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(roomId int, start, end time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
	err := m.DB.QueryRowContext(ctx, stmt,
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
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	stmt := `
		select r.id ,r.room_name from rooms r 
		where r.id not in (
			select rr.room_id from room_restrictions rr 
				where $1 <= rr.end_date 
				and $2 >= rr.start_date
			)
	`
	rows, err := m.DB.QueryContext(ctx, stmt,
		start,
		end,
	)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

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
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var room models.Room
	stmt := `
		select id, room_name, created_at, updated_at from rooms where id = $1
	`

	row := m.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}

// GetUserByID __
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var user models.User

	query := `select id, first_name,last_name, email, password, access_level, created_at, updated_at 
				from users where id = $1`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil

}

// UpdaterUser __
func (m *postgresDBRepo) UpdaterUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update users set first_name = $1,last_name = $2, email = $3, access_level = $4,  updated_at = $5
				where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate receive email and password , then return true if email and password is same to repository
func (m *postgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var id int
	var hashedPassword string

	err := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email).
		Scan(&id, &hashedPassword)
	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return id, hashedPassword, errors.New("wrong password")
	} else if err != nil {
		return id, hashedPassword, err
	}
	return id, hashedPassword, nil

}

// AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var reservations []models.Reservation

	// language=SQL
	query := `select r.id,r.first_name,r.last_name,r.email,r.phone,
				r.start_date,r.end_date,r.created_at,r.updated_at, r.processed, 
				rm.id,rm.room_name
				from reservations r
				left join rooms rm on r.room_id = rm.id
				order by r.start_date
				`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		err = rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Processed,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var reservations []models.Reservation

	// language=SQL
	query := `select r.id,r.first_name,r.last_name,r.email,r.phone,
				r.start_date,r.end_date,r.created_at,r.updated_at, r.processed, 
				rm.id,rm.room_name
				from reservations r
				left join rooms rm on r.room_id = rm.id
				where r.processed = 0
				order by r.start_date
				`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		err = rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Processed,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationByID returns a reservation by id
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var reservation models.Reservation

	query := `select r.id, r.first_name,r.last_name, r.email, r.phone, 
			  r.start_date, r.end_date, r.created_at, r.updated_at, r.processed, r.room_id, rm.room_name
				from reservations r 
				left join rooms rm on r.room_id = rm.id
				where r.id = $1`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)
	if err != nil {
		return reservation, err
	}
	return reservation, nil
}
