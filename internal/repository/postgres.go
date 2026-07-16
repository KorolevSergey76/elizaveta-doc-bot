package repository

import (
	"cosmetologybotliza/internal/domain"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// SaveUser сохраняет или обновляет пользователя
func (r *Repository) SaveUser(id int64, username, firstName string) error {
	query := `INSERT INTO users (telegram_id, username, first_name, role) 
	          VALUES ($1, $2, $3, 'client') 
	          ON CONFLICT (telegram_id) DO NOTHING`
	_, err := r.db.Exec(query, id, username, firstName)
	return err
}

// GetServices возвращает список всех услуг
func (r *Repository) GetServices() ([]domain.Service, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, duration_minutes FROM services")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []domain.Service
	for rows.Next() {
		var s domain.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Price, &s.DurationMinutes); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

// GetMasters возвращает список всех мастеров
func (r *Repository) GetMasters() ([]domain.Master, error) {
	rows, err := r.db.Query("SELECT id, name, bio FROM masters")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var masters []domain.Master
	for rows.Next() {
		var m domain.Master
		if err := rows.Scan(&m.ID, &m.Name, &m.Bio); err != nil {
			return nil, err
		}
		masters = append(masters, m)
	}
	return masters, nil
}
