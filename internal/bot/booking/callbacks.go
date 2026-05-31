package booking

import (
	"cosmetologybotliza/internal/bot/keyboard"
	"cosmetologybotliza/internal/fsm" // Добавили для форматирования строк
	"log"
	"strconv" // Добавили для конвертации ID записи из callback-данных
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) routeCallback(cb *tgbotapi.CallbackQuery) {
	data := cb.Data
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	log.Println("CALLBACK:", data)

	_, _ = h.Bot.Request(tgbotapi.NewCallback(cb.ID, ""))

	switch {

	case data == "noop":
		return

	case data == "book":
		h.StartBooking(cb.Message)
		return

	case data == "my_bookings":
		h.ShowUserBookings(chatID)
		return

	// Одобрение записи админом/мастером
	case strings.HasPrefix(data, "approve_"):
		appointmentIDStr := strings.TrimPrefix(data, "approve_")
		appID, err := strconv.ParseInt(appointmentIDStr, 10, 64)
		if err != nil {
			log.Println("Approve parse ID error:", err)
			return
		}

		// 1. Меняем статус в БД
		err = h.Service.UpdateStatus(appID, "confirmed")
		if err != nil {
			log.Println("Update status error:", err)
			return
		}

		// 2. Достаем запись, чтобы узнать, кому слать уведомление
		app, err := h.Service.GetAppointmentByID(appID)
		if err != nil {
			log.Println("Get appointment error:", err)
			return
		}

		// 3. Обновляем сообщение у админа (убираем кнопки, пишем что подтверждено)
		edit := tgbotapi.NewEditMessageText(chatID, messageID, cb.Message.Text+"\n\n🟢 **Запись успешно подтверждена!**")
		edit.ParseMode = "Markdown"
		_, _ = h.Bot.Request(edit)

		// 4. Отправляем радостную весть клиенту.
		// ВНИМАНИЕ: убедись, что в структуре app у тебя есть Telegram ID клиента (например, ClientTelegramID)
		clientMsg := tgbotapi.NewMessage(app.ClientTelegramID, "🎉 **Ура! Ваша запись подтверждена мастером.**\nС нетерпением ждем вас! 📅")
		clientMsg.ParseMode = "Markdown"
		_, _ = h.Bot.Send(clientMsg)
		return

	// Отклонение записи админом/мастером
	case strings.HasPrefix(data, "reject_"):
		appointmentIDStr := strings.TrimPrefix(data, "reject_")
		appID, err := strconv.ParseInt(appointmentIDStr, 10, 64)
		if err != nil {
			log.Println("Reject parse ID error:", err)
			return
		}

		// 1. Ставим статус cancelled
		_ = h.Service.UpdateStatus(appID, "cancelled")
		app, err := h.Service.GetAppointmentByID(appID)
		if err != nil {
			log.Println("Get appointment error:", err)
			return
		}

		// 2. Обновляем сообщение у админа
		edit := tgbotapi.NewEditMessageText(chatID, messageID, cb.Message.Text+"\n\n🔴 **Запись отклонена.**")
		edit.ParseMode = "Markdown"
		_, _ = h.Bot.Request(edit)

		// 3. Пишем клиенту, что место занято, и предлагаем выбрать другое время
		clientMsg := tgbotapi.NewMessage(app.ClientTelegramID, "😔 **К сожалению, выбранное время уже занято** (например, была запись по телефону).\n\nПожалуйста, выберите другое удобное время!")
		clientMsg.ParseMode = "Markdown"

		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📅 Выбрать другое время", "book"),
			),
		)
		clientMsg.ReplyMarkup = kb
		_, _ = h.Bot.Send(clientMsg)
		return

	// 1.1 Возврат с выбора ДАТЫ назад к выбору МАСТЕРА
	case data == "back_to_masters":
		session, err := h.FSM.Get(chatID)
		if err != nil || session == nil {
			h.Start(cb.Message)
			return
		}

		session.State = fsm.SelectMaster
		_ = h.FSM.Set(chatID, *session)

		kb := keyboard.BuildMasterKeyboard(h.Service)
		kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_services"),
			tgbotapi.NewInlineKeyboardButtonData("📱 В меню", "to_main_menu"),
		))

		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите мастера 👩‍💼")
		edit.ReplyMarkup = &kb
		_, _ = h.Bot.Request(edit)
		return

	// 1.2. Возврат с выбора ВРЕМЕНИ назад к выбору ДАТЫ
	case data == "back_to_dates":
		session, err := h.FSM.Get(chatID)
		if err != nil || session == nil {
			h.Start(cb.Message)
			return
		}

		session.State = fsm.SelectDate
		_ = h.FSM.Set(chatID, *session)

		serviceData, _ := h.Service.GetServiceByID(session.ServiceID)

		kb := keyboard.BuildSmartDateKeyboard(h.Service, session.MasterID, serviceData.DurationMinutes)
		kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_masters"),
			tgbotapi.NewInlineKeyboardButtonData("📱 В меню", "to_main_menu"),
		))

		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите дату 📅")
		edit.ReplyMarkup = &kb
		_, _ = h.Bot.Request(edit)
		return

	// 2.1. Возврат с выбора МАСТЕРА назад к выбору УСЛУГ
	case data == "back_to_services":
		_ = h.FSM.Set(chatID, fsm.Session{
			State:      fsm.SelectService,
			TelegramID: chatID,
		})

		services, err := h.Service.GetServices()
		if err != nil {
			log.Println("error getting services:", err)
			return
		}

		kb := keyboard.BuildServiceKeyboard(services)
		kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 В главное меню", "to_main_menu"),
		))

		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите услугу 💅")
		edit.ReplyMarkup = &kb
		_, _ = h.Bot.Request(edit)
		return

	// Динамические шаги процесса записи
	/*
		главный «маршрутизатор» (роутер) процесса бронирования.
		Именно здесь бот понимает, на каком этапе находится пользователь,
		и передает управление соответствующему методу.
		Когда клиент нажимает любую кнопку, связанную с записью, этот switch читает cb.Data (данные кнопки)
		и запускает нужную логику.
	*/

	//показывает карточку с описанием услуги
	case strings.HasPrefix(data, "service_"):
		h.selectService(cb, data)
	//Бот запоминает выбранного мастера в FSM (сессии)
	case strings.HasPrefix(data, "master_"):
		h.selectMaster(cb, data)
	//Выбор конкретного слота (даты). Бот проверяет доступность времени в базе через h.Service
	case strings.HasPrefix(data, "date_"):
		h.selectDate(cb, data)
	//Выбор конкретного слота (времени). Бот проверяет доступность времени в базе через h.Service
	case strings.HasPrefix(data, "time_"):
		h.selectTime(cb, data)
	//Вызывает метод отмены записи
	case strings.HasPrefix(data, "cancel_"):
		h.cancelBooking(chatID, data, cb)
	//Переход в режим просмотра записей (админская или мастерская часть).
	case data == "master_bookings":
		h.showMasterBookings(chatID, cb.Message.MessageID)
	//Метод handleBack откатывает состояние FSM на один шаг назад (например, от выбора времени к выбору мастера)
	case data == "back":
		h.handleBack(cb)
	//Превращает намерение клиента («Хочу это») в действие («Перехожу к выбору мастера»).
	case strings.HasPrefix(cb.Data, "book_confirm_"):
		h.HandleBookingConfirm(cb)
	//
	default:
		log.Println("unknown callback:", data)
	}
}
