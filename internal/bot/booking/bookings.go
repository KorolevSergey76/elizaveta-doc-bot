package booking

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/auth"
)

// Метод, который позволяет авторизованному персоналу посмотреть список своих записей на текущий момент.
func (h *Handler) showMasterBookings(chatID int64, messageID int) {

	// 1. получаем пользователя
	user, err := h.UserService.GetByTelegramID(chatID)
	if err != nil {
		log.Println(err)

		_, sendErr := h.Bot.Send(
			tgbotapi.NewMessage(chatID, "Ошибка получения пользователя ❌"),
		)
		if sendErr != nil {
			log.Println(sendErr)
		}
		return
	}

	if user == nil {
		_, sendErr := h.Bot.Send(
			tgbotapi.NewMessage(chatID, "Пользователь не найден ❌"),
		)
		if sendErr != nil {
			log.Println(sendErr)
		}
		return
	}

	// 2. проверка роли
	if !auth.IsMaster(user) && !auth.IsAdmin(user) {
		_, sendErr := h.Bot.Send(
			tgbotapi.NewMessage(chatID, "Доступ запрещён ⛔"),
		)
		if sendErr != nil {
			log.Println(sendErr)
		}
		return
	}

	// 3. проверка master_id (нужно только мастеру)
	if auth.IsMaster(user) && user.MasterID == nil {
		_, sendErr := h.Bot.Send(
			tgbotapi.NewMessage(chatID, "Мастер не привязан ❌"),
		)
		if sendErr != nil {
			log.Println(sendErr)
		}
		return
	}

	// 4. определяем masterID
	var masterID int64
	if auth.IsAdmin(user) {
		// админ может смотреть любого мастера
		// пока упрощённо: 1 (можно потом расширить выбором)
		masterID = 1
	} else {
		masterID = *user.MasterID
	}

	// 5. получаем записи
	list, err := h.Service.GetMasterAppointments(masterID)
	if err != nil {
		log.Println(err)

		_, sendErr := h.Bot.Send(
			tgbotapi.NewMessage(chatID, "Ошибка получения записей ❌"),
		)
		if sendErr != nil {
			log.Println(sendErr)
		}
		return
	}

	// 6. если пусто
	if len(list) == 0 {
		_, sendErr := h.Bot.Send(
			tgbotapi.NewMessage(chatID, "Записей нет 📭"),
		)
		if sendErr != nil {
			log.Println(sendErr)
		}
		return
	}

	// 7. формируем ответ
	text := "📋 Записи мастера:\n\n"
	if len(list) == 0 {
		text = "Записей нет 📭"
	} else {
		for _, a := range list {
			text += "🕒 " + a.Time.Format("02.01.2006 15:04") + "\n"
			text += "📌 " + a.Status + "\n"
			if a.Username != "" {
				text += "👤 @" + a.Username + "\n"
			}
			text += "-----------------\n"
		}
	}

	// 8. отправка
	// ИНТЕРАКТИВНОЕ ОБНОВЛЕНИЕ
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	// Добавляем кнопку "Назад" к списку мастеров или главному меню
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "admin_open"),
		),
	)
	edit.ReplyMarkup = &keyboard

	_, _ = h.Bot.Send(edit)
}
