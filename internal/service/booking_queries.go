package service

import "errors"

func (b *BookingService) GetMasterAppointments(
	masterID int64,
) ([]AppointmentView, error) {

	rows, err := b.DB.DB.Query(`
		SELECT
			a.id,
			a.time,
			a.status,
			COALESCE(u.username, '')
		FROM appointments a
		LEFT JOIN users u
			ON u.id = a.user_id
		WHERE a.master_id = $1
		ORDER BY a.time ASC
	`, masterID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []AppointmentView

	for rows.Next() {

		var a AppointmentView

		err := rows.Scan(
			&a.ID,
			&a.Time,
			&a.Status,
			&a.Username,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, a)
	}

	return res, nil
}

func (b *BookingService) GetUserAppointments(
	userID int64,
) ([]AppointmentUserView, error) {

	rows, err := b.DB.DB.Query(`
		SELECT
			a.id,
			a.time,
			a.status,
			s.name,
			m.name
		FROM appointments a
		JOIN services s
			ON s.id = a.service_id
		JOIN masters m
			ON m.id = a.master_id
		WHERE a.user_id = $1
		AND a.status != 'cancelled'
		ORDER BY a.time ASC
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []AppointmentUserView

	for rows.Next() {

		var a AppointmentUserView

		err := rows.Scan(
			&a.ID,
			&a.Time,
			&a.Status,
			&a.ServiceName,
			&a.MasterName,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, a)
	}

	return res, nil
}

func (b *BookingService) CancelAppointment(userID int64, appointmentID int64) error {

	res, err := b.DB.DB.Exec(`
		UPDATE appointments
		SET status = 'cancelled'
		WHERE id = $1
		AND user_id = $2
		AND status != 'cancelled'
	`, appointmentID, userID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}
