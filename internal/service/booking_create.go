package service

import (
	"log"
	"time"
)

func (b *BookingService) Create(
	userID int64,
	serviceID int64,
	masterID int64,
	t time.Time,
) error {

	// timezone салона
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	// нормализуем входящее время в timezone салона
	localTime := t.In(loc)

	// переводим в UTC для хранения в БД
	utcTime := localTime.UTC()

	// получаем duration из services (SNAPSHOT)
	var durationMinutes int

	err = b.DB.DB.QueryRow(`
		SELECT duration_minutes
		FROM services
		WHERE id = $1
	`, serviceID).Scan(&durationMinutes)

	if err != nil {
		return err
	}

	// вставка бронирования
	_, err = b.DB.DB.Exec(`
		INSERT INTO appointments
		(user_id, service_id, master_id, time, duration_minutes, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
	`,
		userID,
		serviceID,
		masterID,
		utcTime,
		durationMinutes,
	)

	if err != nil {
		return err
	}

	log.Println(
		"CREATE APPOINTMENT:",
		"user:", userID,
		"service:", serviceID,
		"master:", masterID,
		"local:", localTime.Format(time.RFC3339),
		"utc:", utcTime.Format(time.RFC3339),
		"duration:", durationMinutes,
	)

	return nil
}
