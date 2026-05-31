package service

import (
	"log"
	"time"
)

type Break struct {
	Start time.Time
	End   time.Time
}

func (b *BookingService) GetBreaks(
	masterID int64,
	date time.Time,
) ([]Break, error) {

	weekday := customWeekday(date)

	rows, err := b.DB.DB.Query(`
        SELECT start_time::text, end_time::text
        FROM master_breaks
        WHERE master_id = $1
        AND weekday = $2
    `, masterID, weekday)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breaks []Break

	for rows.Next() {
		// 1. Считываем время из БД как текст
		var startStr, endStr string
		err := rows.Scan(&startStr, &endStr)
		if err != nil {
			log.Println("Ошибка сканирования времени перерыва:", err)
			return nil, err
		}

		// 2. Парсим строки в объекты time.Time (используем 15:04:05, так как формат в БД - TIME)
		startT, err := time.Parse("15:04:05", startStr)
		if err != nil {
			log.Println("Ошибка парсинга start_time:", err)
			continue
		}
		endT, err := time.Parse("15:04:05", endStr)
		if err != nil {
			log.Println("Ошибка парсинга end_time:", err)
			continue
		}

		// 3. Формируем полную дату перерыва для конкретного дня
		start := time.Date(
			date.Year(), date.Month(), date.Day(),
			startT.Hour(), startT.Minute(), startT.Second(), 0,
			date.Location(),
		)

		end := time.Date(
			date.Year(), date.Month(), date.Day(),
			endT.Hour(), endT.Minute(), endT.Second(), 0,
			date.Location(),
		)

		breaks = append(breaks, Break{
			Start: start,
			End:   end,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return breaks, nil
}
