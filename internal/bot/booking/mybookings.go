package booking

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
Этот код реализует функционал просмотра истории и активных записей пользователя.
По сути, это личный кабинет клиента.
*/

func (h *Handler) ShowUserBookings(chatID int64) {

	user, err := h.UserService.GetByTelegramID(chatID)
	if err != nil || user == nil {
		log.Println(err)
		return
	}

	list, err := h.Service.GetUserAppointments(user.ID)
	if err != nil {

		log.Println(err)

		_, sendErr := h.Bot.Send(
			tgbotapi.NewMessage(
				chatID,
				"Ошибка получения записей ❌",
			),
		)

		if sendErr != nil {
			log.Println(sendErr)
		}

		return
	}

	if len(list) == 0 {

		_, sendErr := h.Bot.Send(
			tgbotapi.NewMessage(
				chatID,
				"У вас нет записей 📭",
			),
		)

		if sendErr != nil {
			log.Println(sendErr)
		}

		return
	}

	_, err = h.Bot.Send(
		tgbotapi.NewMessage(
			chatID,
			"📋 Ваши записи:",
		),
	)

	if err != nil {
		log.Println(err)
	}

	// Загружаем московскую временную зону (UTC+3)
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("Ошибка загрузки локации:", err)
		// Если база/система не знает "Europe/Moscow", используем фиксированное смещение +3 часа
		loc = time.FixedZone("MSK", 3*60*60)
	}

	for _, a := range list {

		var text strings.Builder

		localTime := a.Time.In(loc)

		text.WriteString("📅 ")
		text.WriteString(
			localTime.Format("02.01.2006 15:04"),
		)
		text.WriteString("\n")

		text.WriteString("💅 ")
		text.WriteString(a.ServiceName)
		text.WriteString("\n")

		text.WriteString("👩‍💼 ")
		text.WriteString(a.MasterName)
		text.WriteString("\n")

		text.WriteString("📌 ")
		text.WriteString(a.Status)

		btn := tgbotapi.NewInlineKeyboardMarkup(
			//Первый ряд — кнопка отмены
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"❌ Отменить запись",
					fmt.Sprintf(
						"cancel_%d",
						a.ID,
					),
				),
			),
			// Второй ряд — кнопка возврата
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"⬅️ В главное меню",
					"to_main_menu",
				),
			),
		)

		msg := tgbotapi.NewMessage(
			chatID,
			text.String(),
		)

		msg.ReplyMarkup = btn

		_, err := h.Bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
