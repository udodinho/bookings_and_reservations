package dbrepo

import (
	"context"
	"github.com/udodinho/bookings/internal/models"
	"time"
)

func (m *PostgresDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a new reservation into the database
func (m PostgresDbRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date,
	room_id, created_at, updated_at) 
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
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// InsertRoomRestriction inserts a new room restriction into the database
func (m *PostgresDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restriction (room_id, start_date, end_date, reservation_id, restriction_id, created_at, updated_at) 
		values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.RoomID,
		r.StartDate,
		r.EndDate,
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

// SearchAvailabilityByDatesByRoomID returns true if availability exist for roomID and false if it doesn't.
func (m *PostgresDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `
	select
		 count(id)
	from 
		  room_restriction
	where 
	      room_id = $1
		  $2 < end_date and $3 > start_date ;`

	rows := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := rows.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a list of rooms that are available for the given dates
func (m *PostgresDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		select 
			r.id, r.room_name
		from
				rooms r
		where r.id not in 
		(select rr.room_id from room_restriction rr where $1 < rr.end_date and $2 > rr.start_date)
		`
	rows, err := m.DB.QueryContext(ctx, query, start, end)
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
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

// GetRoomByID returns a room by its ID
func (m *PostgresDbRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select 
			id, room_name, created_at, updated_at from rooms where id = $1
		`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}
