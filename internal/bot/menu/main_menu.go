package menu

import (
	"cosmetologybotliza/internal/auth"
	"cosmetologybotliza/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildMainMenu(user *domain.User) tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	// 1. О мастерах
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👩‍💼 О мастерах", "menu_masters"),
		),
	)

	// 2. Услуги
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💅 Услуги", "menu_services"),
		),
	)

	// 3. Запись
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Записаться", "book"),
		),
	)

	// 4. Мои записи
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои записи", "my_bookings"),
		),
	)

	// 5. Контакты
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📞 Контакты", "menu_contacts"),
		),
	)

	// 6. Как добраться
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📍 Как добраться", "menu_location"),
		),
	)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🎓 Обучение", "menu_education"),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("📢 Перейти на канал", "https://t.me/@koroleva_kosmetolog"),
		tgbotapi.NewInlineKeyboardButtonURL("✉️ Написать врачу", "https://t.me/elizavetadoc02"),
	))

	if auth.IsAdmin(user) {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🛠 Admin", "admin_open"),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
