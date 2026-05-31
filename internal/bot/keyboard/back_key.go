package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func BackButton() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅ Назад", "back"),
		),
	)
}
