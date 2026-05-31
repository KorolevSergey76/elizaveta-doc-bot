package keyboard

import (
	"cosmetologybotliza/internal/domain"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildServiceKeyboard(services []domain.Service) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, s := range services {
		// Форматируем текст кнопки, например: "Маникюр — 2000 ₽"
		buttonText := fmt.Sprintf("%s — %d ₽", s.Name, s.Price)
		callbackData := fmt.Sprintf("service_%d", s.ID)

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
		))
	}

	// Если услуг в базе нет вообще
	if len(rows) == 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Нет доступных услуг 🤷‍♀️", "noop"),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
