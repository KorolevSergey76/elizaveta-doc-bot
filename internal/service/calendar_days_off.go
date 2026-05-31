package service

import "time"

func (b *BookingService) IsDayOff(
	masterID int64,
	date time.Time,
) bool {

	var count int64

	err := b.DB.DB.QueryRow(`
		SELECT COUNT(*)
		FROM master_days_off
		WHERE master_id = $1
		AND off_date = $2
	`,
		masterID,
		date.Format("2006-01-02"),
	).Scan(&count)

	if err != nil {
		return false
	}

	return count > 0
}
