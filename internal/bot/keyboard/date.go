package keyboard

import (
	"fmt"
	"time"

	"cosmetologybotliza/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildSmartDateKeyboard(
	svc service.BookingServiceInterface,
	masterID int64,
	duration int64,
) tgbotapi.InlineKeyboardMarkup {

	today := time.Now()

	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < 7; i++ {

		date := today.AddDate(0, 0, i)

		if svc.IsDayOff(masterID, date) {
			continue
		}

		if !svc.HasFreeSlots(
			masterID,
			date,
			duration,
		) {
			continue
		}

		label := date.Format("02 Jan (Mon)")

		data := fmt.Sprintf(
			"date_%s",
			date.Format("2006-01-02"),
		)

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					label,
					data,
				),
			),
		)
	}

	if len(rows) == 0 {

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Нет свободных дат ❌",
					"noop",
				),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
