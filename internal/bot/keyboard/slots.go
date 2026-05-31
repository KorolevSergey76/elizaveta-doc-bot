package keyboard

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildSlotsKeyboard(
	slots []time.Time,
) tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for i, s := range slots {

		label := s.Format("15:04")

		row = append(row,
			tgbotapi.NewInlineKeyboardButtonData(
				label,
				"time_"+label,
			),
		)

		if len(row) == 2 || i == len(slots)-1 {

			rows = append(rows,
				tgbotapi.NewInlineKeyboardRow(row...),
			)

			row = nil
		}
	}

	if len(rows) == 0 {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Нет слотов", "noop"),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
