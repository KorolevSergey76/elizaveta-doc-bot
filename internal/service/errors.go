package service

import "errors"

var ErrSlotAlreadyBooked = errors.New(
	"slot already booked",
)
