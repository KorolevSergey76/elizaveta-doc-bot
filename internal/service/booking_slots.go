package service

import (
	"log"
	"time"
)

// GetDaySlots генерирует сетку временных слотов для мастера на день по расписанию
func (b *BookingService) GetDaySlots(
	masterID int64,
	date time.Time,
) []time.Time {

	weekday := customWeekday(date)

	var startTime string
	var endTime string

	err := b.DB.DB.QueryRow(`
        SELECT
            start_time::text,
            end_time::text
        FROM master_schedules
        WHERE master_id = $1
        AND weekday = $2
        LIMIT 1
    `, masterID, weekday).Scan(
		&startTime,
		&endTime,
	)

	if err != nil {
		return []time.Time{}
	}

	startParsed, err := time.Parse("15:04:05", startTime)
	if err != nil {
		return []time.Time{}
	}

	endParsed, err := time.Parse("15:04:05", endTime)
	if err != nil {
		return []time.Time{}
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = date.Location()
	}

	start := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		startParsed.Hour(),
		startParsed.Minute(),
		0,
		0,
		loc,
	)

	end := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		endParsed.Hour(),
		endParsed.Minute(),
		0,
		0,
		loc,
	)

	var slots []time.Time

	for t := start; t.Before(end); t = t.Add(30 * time.Minute) {
		slots = append(slots, t)
	}

	return slots
}

// IsSlotFree проверяет по БД, свободен ли конкретный интервал времени у мастера
func (b *BookingService) IsSlotFree(
	masterID int64,
	start time.Time,
	duration int64,
) bool {
	// Вычисляем время окончания новой потенциальной записи
	end := start.Add(time.Duration(duration) * time.Minute)

	var count int64

	// Исправленный и оптимизированный SQL-запрос
	err := b.DB.DB.QueryRow(`
        SELECT COUNT(*)
        FROM appointments a
        JOIN services s ON s.id = a.service_id
        WHERE a.master_id = $1
        AND a.status != 'cancelled'
        AND (
            $2 < (a.time + make_interval(mins => s.duration_minutes)) -- Начало новой записи < Конца существующей
            AND $3 > a.time                                          -- Конец новой записи > Начала существующей
        )
    `, masterID, start, end).Scan(&count)

	if err != nil {
		log.Println("SQL IsSlotFree error:", err)
		return false
	}

	// Если count == 0, значит пересечений нет и слот свободен (true)
	return count == 0
}

// GetFreeSlots возвращает список доступных для записи слотов с учетом занятости и перерывов
func (b *BookingService) GetFreeSlots(
	masterID int64,
	date time.Time,
	duration int64,
) []time.Time {

	if b.IsDayOff(masterID, date) {
		return []time.Time{}
	}

	slots := b.GetDaySlots(masterID, date)
	breaks, _ := b.GetBreaks(masterID, date)

	var free []time.Time
	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)

	for _, s := range slots {

		if s.Before(now) {
			continue
		}

		// ВАЖНОЕ ИСПРАВЛЕНИЕ: Переводим локальное время слота в UTC перед отправкой в IsSlotFree,
		// так как в базе данных время записей хранится в формате UTC.
		slotUTC := s.UTC()

		if b.IsSlotFree(masterID, slotUTC, duration) &&
			!overlapsBreak(s, breaks, duration) {

			free = append(free, s)
		}
	}

	return free
}

// HasFreeSlots проверяет, есть ли вообще свободные слоты на выбранный день
func (b *BookingService) HasFreeSlots(
	masterID int64,
	date time.Time,
	duration int64,
) bool {

	if b.IsDayOff(masterID, date) {
		return false
	}

	slots := b.GetDaySlots(masterID, date)
	breaks, _ := b.GetBreaks(masterID, date)

	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)

	for _, s := range slots {

		if s.Before(now) {
			continue
		}

		// ВАЖНОЕ ИСПРАВЛЕНИЕ: Точно так же приводим к UTC для корректной проверки СУБД
		slotUTC := s.UTC()

		if b.IsSlotFree(masterID, slotUTC, duration) &&
			!overlapsBreak(s, breaks, duration) {

			return true
		}
	}

	return false
}

// CreatePending создает новую запись в БД в статусе 'pending' (ожидает подтверждения)
func (b *BookingService) CreatePending(userID int64, serviceID int64, masterID int64, startTime time.Time) (int64, error) {
	var id int64
	err := b.DB.DB.QueryRow(`
        INSERT INTO appointments (user_id, service_id, master_id, time, status)
        VALUES ($1, $2, $3, $4, 'pending')
        RETURNING id
    `, userID, serviceID, masterID, startTime).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateStatus обновляет статус существующей записи (например: 'confirmed', 'cancelled')
func (b *BookingService) UpdateStatus(appointmentID int64, status string) error {
	res, err := b.DB.DB.Exec(`
        UPDATE appointments 
        SET status = $1 
        WHERE id = $2
    `, status, appointmentID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Запись уже имеет этот статус или не найдена
		return nil
	}

	return nil
}

// GetAppointmentByID находит запись по ее ID и подтягивает Telegram ID клиента
func (b *BookingService) GetAppointmentByID(appointmentID int64) (*Appointment, error) {
	var app Appointment

	err := b.DB.DB.QueryRow(`
        SELECT 
            a.id, 
            a.user_id, 
            a.master_id, 
            a.service_id, 
            a.time, 
            a.status,
            u.telegram_id
        FROM appointments a
        JOIN users u ON u.id = a.user_id
        WHERE a.id = $1
    `, appointmentID).Scan(
		&app.ID,
		&app.UserID,
		&app.MasterID,
		&app.ServiceID,
		&app.Time,
		&app.Status,
		&app.ClientTelegramID,
	)

	if err != nil {
		return nil, err
	}
	return &app, nil
}

// customWeekday переводит стандартный Weekday из Go (0=Вс, 1=Пн...6=Сб)
// в привычный формат ISO (1=Пн, 2=Вт...7=Вс)
func customWeekday(t time.Time) int64 {
	w := int64(t.Weekday())
	if w == 0 {
		return 7
	}
	return w
}

// overlapsBreak проверяет, пересекается ли выбранный слот с перерывами мастера
func overlapsBreak(slot time.Time, breaks []Break, duration int64) bool {
	slotEnd := slot.Add(time.Duration(duration) * time.Minute)

	for _, b := range breaks {
		// Слот пересекается с перерывом, если начало слота раньше конца перерыва
		// И конец слота позже начала перерыва
		if slot.Before(b.End) && slotEnd.After(b.Start) {
			return true
		}
	}

	return false
}
