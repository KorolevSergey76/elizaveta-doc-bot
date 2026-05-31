package domain

import (
	"time"
)

type User struct {
	ID         int64
	TelegramID int64
	Username   string
	FirstName  string
	Role       string
	MasterID   *int64
}

type Service struct {
	ID              int64
	Name            string
	Category        string
	Price           int64
	DurationMinutes int64
	BufferMinutes   int64
	Description     *string
}

type Appointment struct {
	ID               int64
	UserID           int64
	ServiceID        int64
	MasterID         int64
	Time             time.Time
	Status           string
	ClientName       string
	MasterName       string
	ClientTelegramID int64
	ReviewSent       bool
}

type Master struct {
	ID   int64
	Name string
}
