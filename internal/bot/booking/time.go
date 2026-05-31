package booking

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) selectTime(cb *tgbotapi.CallbackQuery, data string) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	session, err := h.FSM.Get(chatID)
	if err != nil {
		log.Println("FSM Error:", err)
		return
	}

	timePart := strings.TrimPrefix(data, "time_")
	parsed, err := time.Parse("15:04", timePart)
	if err != nil {
		log.Println("Parse time error:", err)
		return
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("Load location error:", err)
		return
	}

	localDateTime := time.Date(
		session.Date.Year(), session.Date.Month(), session.Date.Day(),
		parsed.Hour(), parsed.Minute(), 0, 0, loc,
	)
	t := localDateTime.UTC()

	serviceData, err := h.Service.GetServiceByID(session.ServiceID)
	if err != nil {
		log.Println("Get service error:", err)
		return
	}

	// Проверяем занятость
	if !h.Service.IsSlotFree(session.MasterID, t, serviceData.DurationMinutes) {
		alert := tgbotapi.NewCallbackWithAlert(cb.ID, "Упс! Это время уже занято. Пожалуйста, выберите другое! ❌")
		_, _ = h.Bot.Request(alert)
		return
	}

	user, err := h.UserService.GetByTelegramID(chatID)
	if err != nil {
		log.Println("Get user error:", err)
		return
	}

	// 1. Создаем запись со статусом 'pending' (это нужно поправить в твоем методе Service.Create)
	// Предположим, метод возвращает ID созданной записи (appointmentID), он нам нужен для кнопок!
	appointmentID, err := h.Service.CreatePending(user.ID, session.ServiceID, session.MasterID, t)
	if err != nil {
		log.Println("Create pending appointment error:", err)
		return
	}

	// 2. Отвечаем клиенту, что заявка на рассмотрении
	kbAfterBooking := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 В главное menu", "to_main_menu"),
		),
	)
	editOldMsg := tgbotapi.NewEditMessageText(chatID, messageID,
		"✨ Ваша заявка отправлена на подтверждение! \nМастер проверит график и вам придет уведомление. Скрестили пальчики 🤞")
	editOldMsg.ReplyMarkup = &kbAfterBooking
	_, _ = h.Bot.Request(editOldMsg)

	_ = h.FSM.Clear(chatID)
	_, _ = h.Bot.Request(tgbotapi.NewCallback(cb.ID, "Заявка отправлена"))

	// 3. ОТПРАВЛЯЕМ УВЕДОМЛЕНИЕ МАСТЕРУ ИЛИ АДМИНУ
	// Для примера шлем админу (Sergey, твой ID 583008737).
	// В идеале тут нужно получать telegram_id мастера из таблицы users по session.MasterID
	adminTelegramID := int64(583008737)

	notificationText := fmt.Sprintf(
		"🔔 **Новая заявка на запись!**\n\n"+
			"👤 Клиент: %s (@%s)\n"+
			"💅 Услуга: %s\n"+
			"📅 Дата: %s\n"+
			"⏰ Время: %s",
		user.FirstName, user.Username, serviceData.Name,
		localDateTime.Format("02.01.2006"), localDateTime.Format("15:04"),
	)

	// Кнопки для админа/мастера. Передаем ID записи прямо в callback_data
	kbAdmin := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", fmt.Sprintf("approve_%d", appointmentID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отклонить", fmt.Sprintf("reject_%d", appointmentID)),
		),
	)

	msgToAdmin := tgbotapi.NewMessage(adminTelegramID, notificationText)
	msgToAdmin.ParseMode = "Markdown"
	msgToAdmin.ReplyMarkup = kbAdmin
	_, _ = h.Bot.Send(msgToAdmin)
}
