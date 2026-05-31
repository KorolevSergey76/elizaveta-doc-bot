package service

import (
	"cosmetologybotliza/internal/storage"
	"log"
)

type UserService struct {
	DB *storage.Postgres
}

func (u *UserService) SaveUser(
	telegramID int64,
	username string,
	firstName string,
) error {

	_, err := u.DB.DB.Exec(`
		INSERT INTO users (
			telegram_id,
			username,
			first_name,
			role
		)
		VALUES ($1, $2, $3, 'client')
		ON CONFLICT (telegram_id)
		DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name
	`,
		telegramID,
		username,
		firstName,
	)

	return err
}

func (s *UserService) GetTelegramIDByMasterID(masterID int64) int64 {
	var telegramID int64

	// Предполагаю, что у вас есть таблица masters или users,
	// где хранится связь master_id и telegram_id
	err := s.DB.DB.QueryRow(`
        SELECT telegram_id 
        FROM users 
        WHERE master_id = $1
    `, masterID).Scan(&telegramID)

	if err != nil {
		log.Printf("Ошибка при поиске telegram_id мастера %d: %v", masterID, err)
		return 0
	}

	return telegramID
}
