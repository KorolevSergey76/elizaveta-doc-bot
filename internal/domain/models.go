package domain

// User представляет пользователя бота для учета в базе данных.
type User struct {
	ID         int64
	TelegramID int64
	Username   string
	FirstName  string
	Role       string // "client" или "admin"
}

// Service представляет косметическую услугу.
type Service struct {
	ID              int64
	Name            string
	Description     *string
	Price           int
	DurationMinutes int
}

// Master представляет мастера.
type Master struct {
	ID   int64
	Name string
	Bio  string
}
