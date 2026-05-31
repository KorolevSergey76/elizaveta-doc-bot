package fsm

import (
	"context"
	"cosmetologybotliza/internal/storage"
	"time"
)

type State string

const (
	None          State = "none"
	SelectService State = "service"
	SelectMaster  State = "master"
	SelectDate    State = "date"
	SelectTime    State = "time"
)

type Session struct {
	State      State `json:"state"`
	TelegramID int64
	ServiceID  int64     `json:"service_id"`
	MasterID   int64     `json:"master_id"`
	Date       time.Time `json:"date"`
	PrevState  State     `json:"prev_state"`
}

// Структура FSM определена здесь ОДИН раз
type FSM struct {
	Redis *storage.RedisClient
	Ctx   context.Context
}
