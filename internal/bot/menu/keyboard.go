package menu

import (
	"cosmetologybotliza/internal/domain"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MainMenuKeyboard возвращает статичное главное меню
func MainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👩‍⚕️ Информация о враче", "menu_doctor"),
			tgbotapi.NewInlineKeyboardButtonData("📋 Услуги", "menu_services"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🗓 Записаться", "https://dikidi.net/ru/profile/kosmetolog_korolyova_elizaveta_503969"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎓 Обучение косметологов", "menu_education"),
			tgbotapi.NewInlineKeyboardButtonData("⭐ Отзывы", "menu_reviews"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📍 Как добраться", "menu_location"),
			tgbotapi.NewInlineKeyboardButtonData("📞 Контакты", "menu_contacts"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📢 Перейти на канал", "https://t.me/@koroleva_kosmetolog"),
			tgbotapi.NewInlineKeyboardButtonURL("💬 Написать врачу", "https://t.me/@elizavetadoc02"),
		),
	)
}

// ServiceKeyboard создает список услуг с URL-кнопками для записи
func ServiceKeyboard(services []domain.Service, dikidiURL string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, s := range services {
		buttonText := fmt.Sprintf("%s — %d ₽", s.Name, s.Price)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(buttonText, dikidiURL),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_main"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// MasterKeyboard создает список мастеров
func MasterKeyboard(masters []domain.Master, dikidiURL string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, m := range masters {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(m.Name, dikidiURL),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_main"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// BackButton возвращает клавиатуру с одной кнопкой "Назад"
func BackButton() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_main"),
		),
	)
}

// DoctorListKeyboard показывает список врачей
func DoctorListKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👩‍⚕️Елизавета Королева", "doctor_eliza"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👩‍⚕️Лапина Полина", "doctor_polina"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_main"),
		),
	)
}

func ServicesCategoriesKeyboard() tgbotapi.InlineKeyboardMarkup {
	services := []struct {
		Text string
		Data string
	}{
		{"👄 Увеличение губ", "service_lips_up"},
		{"🧪 Выведение губ", "service_lips_down"},
		{"💉 Ботулинотерапия", "service_botox"},
		{"👤 Контурная пластика", "service_contour"},
		{"💧 Биоревитализация/Мезо", "service_bio"},
		{"⚙️ Аппаратная космет.", "service_hardware"},
		{"🗣 Консультация", "service_consultation"},
		{"🧴 Уходовые программы", "service_care"},
		{"🧪 Пилинги", "service_peeling"},
		{"🧼 Чистка лица", "service_cleaning"},
		{"🪒 Депиляция", "service_depilation"},
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(services); i += 2 {
		row := tgbotapi.NewInlineKeyboardRow()
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(services[i].Text, services[i].Data))
		if i+1 < len(services) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(services[i+1].Text, services[i+1].Data))
		}
		rows = append(rows, row)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_main"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func ContactsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📞 Позвонить", "show_phone"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("VK", "https://vk.com/lizaganicheva"),
			tgbotapi.NewInlineKeyboardButtonURL("Instagram", "https://www.instagram.com/dr.koroleva_elizaveta"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_main"),
		),
	)
}

func LipsQuizKeyboard(step string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	switch step {
	case "q1":
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "quiz_q1_yes"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "quiz_q1_no"),
		))
	case "q2":
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "quiz_q2_yes"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "quiz_q2_no"),
			tgbotapi.NewInlineKeyboardButtonData("Не знаю", "quiz_q2_idk"),
		))
	case "q3":
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "quiz_q3_yes"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "quiz_q3_no"),
			tgbotapi.NewInlineKeyboardButtonData("Не помню", "quiz_q3_idk"),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_services")))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
