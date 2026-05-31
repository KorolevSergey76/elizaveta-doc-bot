package service

import (
	"time"

	"cosmetologybotliza/internal/domain"
)

type BookingServiceInterface interface {
	GetMasters() ([]Master, error)

	CreatePending(userID int64, serviceID int64, masterID int64, startTime time.Time) (int64, error)

	UpdateStatus(appointmentID int64, status string) error
	GetAppointmentByID(appointmentID int64) (*Appointment, error)

	GetServiceByID(id int64) (*domain.Service, error)

	GetFreeSlots(
		masterID int64,
		date time.Time,
		duration int64,
	) []time.Time

	IsSlotFree(
		masterID int64,
		start time.Time,
		duration int64,
	) bool

	IsDayOff(
		masterID int64,
		date time.Time,
	) bool

	Create(
		userID int64,
		serviceID int64,
		masterID int64,
		t time.Time,
	) error

	HasFreeSlots(
		masterID int64,
		date time.Time,
		duration int64,
	) bool

	GetMasterAppointments(masterID int64) ([]AppointmentView, error)

	GetUserAppointments(userID int64) ([]AppointmentUserView, error)

	CancelAppointment(userID int64, appointmentID int64) error

	GetBreaks(
		masterID int64,
		date time.Time,
	) ([]Break, error)

	GetServices() ([]domain.Service, error)

	GetFinishedAppointments(threshold time.Time) ([]domain.Appointment, error)
	MarkReviewSent(appointmentID int64) error

	DeleteOldCancelledAppointments(daysOld int) (int64, error)
}
