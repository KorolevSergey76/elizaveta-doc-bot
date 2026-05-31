package fsm

import (
	"context"
	"cosmetologybotliza/internal/storage"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Конструктор
func NewFSM(client *storage.RedisClient) *FSM {
	return &FSM{
		Redis: client,
		Ctx:   context.Background(),
	}
}

func key(userID int64) string {
	return "cosmo:v1:fsm:" + strconv.FormatInt(userID, 10)
}

func (f *FSM) Set(userID int64, session Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return f.Redis.Client.Set(f.Ctx, key(userID), data, time.Minute*30).Err()
}

func (f *FSM) Get(userID int64) (*Session, error) {
	val, err := f.Redis.Client.Get(f.Ctx, key(userID)).Result()
	if err == redis.Nil {
		return &Session{State: None}, nil
	}
	if err != nil {
		return nil, err
	}
	var session Session
	err = json.Unmarshal([]byte(val), &session)
	return &session, err
}

func (f *FSM) Clear(userID int64) error {
	return f.Redis.Client.Del(f.Ctx, key(userID)).Err()
}

func (f *FSM) Move(userID int64, next State, session *Session) error {
	session.PrevState = session.State
	session.State = next
	return f.Set(userID, *session)
}
