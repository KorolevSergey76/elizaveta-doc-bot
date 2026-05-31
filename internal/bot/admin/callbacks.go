package admin

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/auth"
)

/*
Этот файл — сердце административной части бота.
Он отвечает за то, чтобы администратор мог просматривать записи к конкретным мастерам,
не обращаясь напрямую к базе данных.
Перехватывает все нажатия кнопок, которые начинаются с префикса admin_,
и решает, что именно должен увидеть администратор
*/

// Диспетчер админки. Это точка входа для всех нажатий кнопок в админ-меню.
func (h *Handler) HandleCallback(cb *tgbotapi.CallbackQuery) {

	if strings.HasPrefix(cb.Data, "admin_approve_") {
		// Вырезаем ID записи
		idStr := cb.Data[len("admin_approve_"):]
		appointmentID, _ := strconv.ParseInt(idStr, 10, 64)

		// Обновляем статус в базе данных
		err := h.Service.UpdateStatus(appointmentID, "confirmed")
		if err != nil {
			h.Bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "Ошибка обновления статуса ❌"))
			return
		}

		// Редактируем сообщение мастера, чтобы он видел, что запись принята
		edit := tgbotapi.NewEditMessageText(cb.Message.Chat.ID, cb.Message.MessageID,
			cb.Message.Text+"\n\n✅ **Запись подтверждена!**")
		h.Bot.Send(edit)

		// Убираем кнопку
		h.Bot.Request(tgbotapi.NewCallback(cb.ID, "Успешно!"))
	}

	if cb.Message == nil {
		return
	}
	callbackCfg := tgbotapi.NewCallback(cb.ID, "Обновляю...")
	h.Bot.Request(callbackCfg)

	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	_, _ = h.Bot.Request(tgbotapi.NewCallback(cb.ID, ""))

	user, err := h.User.GetByTelegramID(chatID)
	if err != nil {
		log.Println(err)
		return
	}

	if !auth.IsAdmin(user) {
		_, _ = h.Bot.Send(tgbotapi.NewMessage(chatID, "⛔ Нет доступа"))
		return
	}

	switch {

	// открыть admin panel
	case cb.Data == "admin_open":
		h.OpenAdminMenu(chatID, cb.Message.MessageID)

	// список мастеров
	case cb.Data == "admin_masters":
		h.showMasters(chatID, cb.Message.MessageID)

	case cb.Data == "admin_my_bookings":
		// Вызываем метод для просмотра своих записей (тот, что мы обсуждали)
		h.showMasterBookings(chatID, messageID, cb.Data)
	case cb.Data == "menu_main":
		// Если кнопка ведет в меню, мы вызываем метод из пакета menu
		// Либо, если вы хотите просто удалить админское сообщение:
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
		_, err = h.Bot.Send(deleteMsg)
		if err != nil {
			log.Println(err)
			return
		}

	// записи мастера
	case len(cb.Data) > 8 && cb.Data[:8] == "admin_m_":
		h.showMasterBookings(chatID, cb.Message.MessageID, cb.Data)

	}
}

// Список мастеров
func (h *Handler) showMasters(chatID int64, messageID int) {

	//Запрашивает всех мастеров из БД через h.Service.GetMasters()
	masters, err := h.Service.GetMasters()
	if err != nil {
		log.Println(err)
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, m := range masters {

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(m.Name, fmt.Sprintf("admin_m_%d", m.ID)),
		),
		)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "admin_open")),
	)

	text := "👩‍💼 Выберите мастера"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	if messageID > 0 {
		// Редактируем сообщение, если оно есть
		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ReplyMarkup = &keyboard
		_, _ = h.Bot.Send(edit)
	} else {
		// Или шлем новое
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		_, _ = h.Bot.Send(msg)
	}
}

// Просмотр записей. Здесь происходит получение и форматирование данных
func (h *Handler) showMasterBookings(chatID int64, messageID int, data string) {
	var masterID int64
	var err error

	// 1. ОПРЕДЕЛЯЕМ ID МАСТЕРА
	if data == "admin_my_bookings" {
		// Если это личный кабинет, берем ID из профиля пользователя
		user, err := h.User.GetByTelegramID(chatID)
		if err != nil || user == nil || user.MasterID == nil {
			h.Bot.Send(tgbotapi.NewMessage(chatID, "Не удалось найти ваш профиль мастера ❌"))
			return
		}
		masterID = *user.MasterID
	} else if strings.HasPrefix(data, "admin_m_") {
		// Если это выбор из списка мастеров, парсим ID
		idStr := data[len("admin_m_"):]
		mID, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("Ошибка парсинга ID:", err)
			return
		}
		masterID = int64(mID)
	} else {
		log.Println("Неизвестный формат данных:", data)
		return
	}

	// 2. ПОЛУЧАЕМ ЗАПИСИ (код без изменений)
	list, err := h.Service.GetMasterAppointments(masterID)
	if err != nil {
		log.Println("Ошибка получения записей:", err)
		_, _ = h.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения записей ❌"))
		return
	}

	// 3. ФОРМИРУЕМ ТЕКСТ (код без изменений)
	text := "📋 Записи мастера:\n\n"
	if len(list) == 0 {
		text += "Записей пока нет 📭"
	} else {
		loc, _ := time.LoadLocation("Europe/Moscow")
		for _, a := range list {
			localTime := a.Time.In(loc)
			text += "🕒 " + localTime.Format("02.01.2006 15:04") + "\n"
			if a.Status == "cancelled" {
				text += "❌ *Status: " + a.Status + "*\n"
			} else if a.Status == "confirmed" {
				text += "✅ *Status: " + a.Status + "*\n"
			} else {
				text += "ℹ️ *Status: " + a.Status + "*\n"
			}
			if a.Username != "" {
				text += "👤 @" + a.Username + "\n"
			}
			text += "-----------------\n"
		}
	}

	// 4. КЛАВИАТУРА (динамический возврат назад)
	// Если пришли из "Мои записи", кнопка "Назад" должна вести в админ-меню,
	// если из списка мастеров - в список мастеров.
	backData := "admin_masters"
	if data == "admin_my_bookings" {
		backData = "admin_open"
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", data), // Просто повторяем текущий запрос
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", backData),
		),
	)

	// 5. ОБНОВЛЕНИЕ СООБЩЕНИЯ
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ReplyMarkup = &keyboard
	edit.ParseMode = "Markdown"

	_, err = h.Bot.Send(edit)
	if err != nil {
		// Игнорируем ошибку, если сообщение не изменилось
		if !strings.Contains(err.Error(), "message is not modified") {
			log.Printf("Ошибка при обновлении сообщения: %v", err)
		}
	}
}

func (h *Handler) NotifyMaster(masterTelegramID int64, appointmentID int64, details string) {
	// Создаем кнопку для мгновенного подтверждения
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить запись", fmt.Sprintf("admin_approve_%d", appointmentID)),
		),
	)

	msg := tgbotapi.NewMessage(masterTelegramID, "🔔 **Новая запись!**\n\n"+details)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "Markdown"

	_, _ = h.Bot.Send(msg)
}
