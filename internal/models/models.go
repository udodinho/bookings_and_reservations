package models

import "time"

// Reservation is a struct that represents a reservation
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	RoomID    int
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// User is a struct that represents a user
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is a struct that represents a room
type Room struct {
	ID        int
	RooName   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoomRestriction is a struct that represents a room restriction
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

// Restriction is a struct that represents a restriction
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
