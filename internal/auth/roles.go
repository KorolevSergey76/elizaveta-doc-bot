package auth

import "cosmetologybotliza/internal/domain"

/*
Этот пакет auth (аутентификация/авторизация) отвечает за проверку прав доступа пользователя в боте.
Он определяет, кто есть кто, и позволяет ограничивать доступ к функциям
(например, чтобы обычный клиент не мог зайти в меню администратора).
*/

type Role string

/*
Это «словарь» ролей. Вместо того чтобы писать в коде строковые значения "admin" или "client"
(где можно легко допустить опечатку), вы используете константы. Это делает код надежным.
*/
const (
	RoleClient Role = "client"
	RoleMaster Role = "master"
	RoleAdmin  Role = "admin"
)

// Эти функции позволяют в любом месте кода легко проверить права.
func IsClient(u *domain.User) bool {
	return u != nil && Role(u.Role) == RoleClient
}

func IsMaster(u *domain.User) bool {
	return u != nil && Role(u.Role) == RoleMaster
}

func IsAdmin(u *domain.User) bool {
	return u != nil && Role(u.Role) == RoleAdmin
}

// гибкая проверка. Можно проверить, имеет ли пользователь хотя бы одну из указанных ролей.
func RequireRole(u *domain.User, roles ...Role) bool {
	if u == nil {
		return false
	}

	userRole := Role(u.Role)

	for _, r := range roles {
		if userRole == r {
			return true
		}
	}

	return false
}
