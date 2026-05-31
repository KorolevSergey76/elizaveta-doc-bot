package keyboard

import (
	"fmt"
	"log"

	"cosmetologybotliza/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildMasterKeyboard(
	svc service.BookingServiceInterface,
) tgbotapi.InlineKeyboardMarkup {

	masters, err := svc.GetMasters()
	if err != nil {
		log.Println("failed to load masters:", err)

		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Ошибка загрузки",
					"noop",
				),
			),
		)
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, m := range masters {

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					m.Name,
					fmt.Sprintf("master_%d", m.ID),
				),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
