package service

import (
	"cosmetologybotliza/internal/storage"
	"time"
)

type BookingService struct {
	DB *storage.Postgres
}

type Master struct {
	ID   int64
	Name string
}

type AppointmentView struct {
	ID       int64
	Time     time.Time
	Status   string
	Username string
}

type AppointmentUserView struct {
	ID          int64
	Time        time.Time
	Status      string
	ServiceName string
	MasterName  string
}

var _ BookingServiceInterface = (*BookingService)(nil)

type Appointment struct {
	ID               int64
	UserID           int64
	MasterID         int64
	ServiceID        int64
	Time             time.Time
	Status           string
	ClientTelegramID int64 // Это поле мы считываем в SQL-запросе!
}
