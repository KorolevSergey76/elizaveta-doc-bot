package service

import (
	"database/sql"

	"cosmetologybotliza/internal/domain"
)

func (u *UserService) GetByTelegramID(telegramID int64) (*domain.User, error) {

	var user domain.User

	err := u.DB.DB.QueryRow(`
		SELECT
			id,
			telegram_id,
			username,
			first_name,
			role,
			master_id
		FROM users
		WHERE telegram_id = $1
		LIMIT 1
	`, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.Role,
		&user.MasterID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
