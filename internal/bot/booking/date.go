package booking

import (
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/bot/keyboard"
	"cosmetologybotliza/internal/fsm"
)

/*
Этот код — логический узел выбора времени в процессе бронирования.
Его главная задача — «отфильтровать» доступные временные слоты и показать пользователю только те,
на которые реально можно записаться.
*/

func (h *Handler) selectDate(
	cb *tgbotapi.CallbackQuery,
	data string,
) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	session, err := h.FSM.Get(chatID)
	if err != nil {
		log.Println(err)
		return
	}

	dateStr := strings.TrimPrefix(data, "date_")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Println(err)
		return
	}

	serviceData, err := h.Service.GetServiceByID(session.ServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	session.State = fsm.SelectTime
	session.Date = date

	_ = h.FSM.Set(chatID, *session)

	// 1. Получаем все сгенерированные слоты времени
	allSlots := h.Service.GetFreeSlots(
		session.MasterID,
		date,
		serviceData.DurationMinutes,
	)

	// ФИЛЬТРУЕМ ТОЛЬКО РЕАЛЬНО СВОБОДНЫЕ СЛОТЫ
	var freeSlots []time.Time

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println(err)
		return
	}

	for _, slotTime := range allSlots {
		// Собираем полную дату и время для проверки в базе данных
		fullDateTime := time.Date(
			date.Year(),
			date.Month(),
			date.Day(),
			slotTime.Hour(),
			slotTime.Minute(),
			0,
			0,
			loc,
		)

		// Проверяем через твой сервис: если слот свободен, добавляем его в список
		if h.Service.IsSlotFree(session.MasterID, fullDateTime, serviceData.DurationMinutes) {
			freeSlots = append(freeSlots, slotTime)
		}
	}

	// 2. Строим клавиатуру передавая ей только РЕАЛЬНО свободные времена
	kb := keyboard.BuildSlotsKeyboard(freeSlots)

	// Прикрепляем кнопки навигации внизу
	kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_dates"),
		tgbotapi.NewInlineKeyboardButtonData("📱 В меню", "to_main_menu"),
	))

	edit := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		"Выберите время ⏰",
	)
	edit.ReplyMarkup = &kb

	_, err = h.Bot.Request(edit)
	if err != nil {
		log.Println("edit error:", err)
	}
}
