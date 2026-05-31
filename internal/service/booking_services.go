package service

import (
	"cosmetologybotliza/internal/domain"
	"log"
	"time"
)

func (b *BookingService) GetMasters() ([]Master, error) {

	rows, err := b.DB.DB.Query(`
		SELECT id, name
		FROM masters
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var masters []Master

	for rows.Next() {

		var m Master

		err := rows.Scan(
			&m.ID,
			&m.Name,
		)

		if err != nil {
			return nil, err
		}

		masters = append(masters, m)
	}

	return masters, nil
}

func (b *BookingService) GetServiceByID(
	id int64,
) (*domain.Service, error) {

	var s domain.Service

	err := b.DB.DB.QueryRow(`
		SELECT
			id,
			name,
			price,
			duration_minutes,
			buffer_minutes
		FROM services
		WHERE id = $1
		LIMIT 1
	`, id).Scan(
		&s.ID,
		&s.Name,
		&s.Price,
		&s.DurationMinutes,
		&s.BufferMinutes,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (b *BookingService) GetUserID(telegramID int64) (int64, error) {
	var id int64

	err := b.DB.DB.QueryRow(`
		SELECT id FROM users WHERE telegram_id = $1
	`, telegramID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (b *BookingService) GetServices() ([]domain.Service, error) {
	rows, err := b.DB.DB.Query(`
        SELECT id, name, price, duration_minutes, buffer_minutes
        FROM services
        ORDER BY id ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []domain.Service
	for rows.Next() {
		var s domain.Service
		err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.DurationMinutes, &s.BufferMinutes)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

// GetFinishedAppointments находит записи, время которых истекло и отзыв не отправлен
func (b *BookingService) GetFinishedAppointments(threshold time.Time) ([]domain.Appointment, error) {
	query := `
        SELECT a.id, a.user_id, a.master_id, u.first_name, m.name, u.telegram_id, a.time 
        FROM appointments a
        JOIN users u ON a.user_id = u.id
        JOIN masters m ON a.master_id = m.id
        WHERE (a.time + (a.duration_minutes * INTERVAL '1 minute')) < $1 
        AND a.review_sent = false`

	rows, err := b.DB.DB.Query(query, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Appointment
	for rows.Next() {
		var a domain.Appointment
		// Обратите внимание: поля ClientName и т.д. должны быть в domain.Appointment
		err := rows.Scan(&a.ID, &a.UserID, &a.MasterID, &a.ClientName, &a.MasterName, &a.ClientTelegramID, &a.Time)
		if err != nil {
			log.Println("Ошибка сканирования:", err)
			continue
		}
		list = append(list, a)
	}
	return list, nil
}

// MarkReviewSent отмечает запись как "отзыв запрошен"
func (b *BookingService) MarkReviewSent(appointmentID int64) error {
	_, err := b.DB.DB.Exec(`UPDATE appointments SET review_sent = true WHERE id = $1`, appointmentID)
	return err
}

func (b *BookingService) DeleteOldCancelledAppointments(daysOld int) (int64, error) {
	// Удаляем все записи со статусом 'cancelled', которые были созданы/обновлены более X дней назад
	res, err := b.DB.DB.Exec(`
        DELETE FROM appointments 
        WHERE status = 'cancelled' 
        AND time < NOW() - ($1 * INTERVAL '1 day')
    `, daysOld)

	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
