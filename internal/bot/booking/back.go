package booking

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/bot/keyboard"
	"cosmetologybotliza/internal/fsm"
)

/*
Этот код — движок навигации "Назад" для процесса записи клиента к мастеру.
В ботах, где используется FSM (Finite State Machine — Конечный автомат), важно не просто «открыть предыдущее меню»,
а еще и откатить состояние сессии пользователя, чтобы бот «забыл» выбор, сделанный на текущем шаге.
Как работает этот механизм:
Запрос сессии: Бот берет текущий статус пользователя (session.State) из Redis через h.FSM.Get(chatID).
Это позволяет понять, на каком этапе записи находится клиент.
Смена состояния: В switch для каждого шага вы вручную меняете состояние на «шаг назад» (например, если пользователь
был на SelectDate, вы меняете его состояние на SelectMaster).
Перерисовка интерфейса: Бот выполняет обратный путь — заново делает запрос в базу данных или сервис,
чтобы получить список мастеров/услуг, и присылает обновленную клавиатуру через EditMessageText.
Сохранение: В самом конце вызывается h.FSM.Set, чтобы записать обновленное состояние сессии обратно в базу (Redis).
*/

func (h *Handler) handleBack(cb *tgbotapi.CallbackQuery) {

	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	session, err := h.FSM.Get(chatID)
	if err != nil {
		log.Println(err)
		return
	}

	switch session.State {

	/*Шаг назад к услугам: Пользователь хочет сменить услугу.
	Вы запрашиваете список услуг из БД (h.Service.GetServices()) и перерисовываете меню услуг.*/
	case fsm.SelectMaster:
		session.State = fsm.SelectService

		// 1. Получаем актуальные услуги из базы данных
		services, err := h.Service.GetServices()
		if err != nil {
			log.Println("back error (get services):", err)
			// Если база недоступна, отправляем пользователю уведомление, чтобы бот не завис
			_, err = h.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при загрузке услуг ❌"))
			if err != nil {
				log.Println(err)
			}
			return
		}

		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите услугу 💅")

		// 2. Передаем полученные услуги в конструктор клавиатуры
		kb := keyboard.BuildServiceKeyboard(services)
		edit.ReplyMarkup = &kb

		_, err = h.Bot.Request(edit)
		if err != nil {
			log.Println("edit error:", err)
		}

		/*Шаг назад к мастерам: Пользователь передумал выбирать дату и хочет посмотреть другого мастера.
		Бот перерисовывает клавиатуру мастеров.*/
	case fsm.SelectDate:
		session.State = fsm.SelectMaster

		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите мастера 👩‍💼")

		kb := keyboard.BuildMasterKeyboard(h.Service)
		edit.ReplyMarkup = &kb

		_, err = h.Bot.Request(edit)
		if err != nil {
			log.Println("edit error:", err)
		}

	/*Шаг назад к датам: Самый сложный случай. Нужно заново посчитать свободные слоты (GetFreeSlots),
	так как время (слоты) жестко привязано к конкретной дате и выбранному мастеру.*/
	case fsm.SelectTime:
		session.State = fsm.SelectDate

		serviceData, _ := h.Service.GetServiceByID(session.ServiceID)

		slots := h.Service.GetFreeSlots(
			session.MasterID,
			session.Date,
			serviceData.DurationMinutes,
		)

		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите дату 📅")

		kb := keyboard.BuildSlotsKeyboard(slots)
		edit.ReplyMarkup = &kb

		_, err = h.Bot.Request(edit)
		if err != nil {
			log.Println("edit error:", err)
		}
	}

	_ = h.FSM.Set(chatID, *session)
}
