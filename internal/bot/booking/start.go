package booking

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/bot/keyboard"
	"cosmetologybotliza/internal/bot/menu"
	"cosmetologybotliza/internal/domain"
	"cosmetologybotliza/internal/fsm"
	"cosmetologybotliza/internal/service"
)

/*
Этот фрагмент кода содержит логику инициализации пользователя и запуска процесса записи в вашем боте.
Это "фундамент" взаимодействия: с чего всё начинается, когда пользователь нажимает /start или кнопку записи.
*/

type Handler struct {
	Bot         *tgbotapi.BotAPI
	Service     service.BookingServiceInterface
	UserService *service.UserService
	FSM         *fsm.FSM
}

func (h *Handler) Start(msg *tgbotapi.Message) {

	_ = h.UserService.SaveUser(
		msg.From.ID,
		msg.From.UserName,
		msg.From.FirstName,
	)

	user, _ := h.UserService.GetByTelegramID(msg.From.ID)

	if user == nil {
		user = &domain.User{Role: "client"}
	}

	// 2. Настраиваем ReplyKeyboard (ту, что внизу под строкой ввода)
	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📱 Главное меню"),
		),
	)
	replyMenu.ResizeKeyboard = true // Делает кнопку компактной, а не огромной

	// 3. Собираем текст и Inline-клавиатуру
	text := "Добро пожаловать! Используйте меню ниже для навигации ✨\n\n💅 Главное меню"

	// Используем вашу функцию из пакета menu
	inlineMenu := menu.BuildMainMenu(user)

	// Отправляем приветственное сообщение, к которому привязана нижняя кнопка
	welcomeMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
	welcomeMsg.ReplyMarkup = inlineMenu // Inline-клавиатура (кнопки с действиями)

	// Отправляем нижнюю кнопку
	msgReply := tgbotapi.NewMessage(msg.Chat.ID, "Вы находитесь в главном меню")
	msgReply.ReplyMarkup = replyMenu
	h.Bot.Send(msgReply)

	// И сразу же сообщение с самим меню
	msgInline := tgbotapi.NewMessage(msg.Chat.ID, "Выберите действие:")
	msgInline.ReplyMarkup = menu.BuildMainMenu(user)
	h.Bot.Send(msgInline)

	if _, err := h.Bot.Send(welcomeMsg); err != nil {
		log.Println("Ошибка отправки:", err)
	}
}

func (h *Handler) StartBooking(msg *tgbotapi.Message) {
	err := h.FSM.Set(msg.Chat.ID, fsm.Session{
		State:      fsm.SelectService,
		TelegramID: msg.Chat.ID,
	})
	if err != nil {
		log.Println(err)
		return
	}

	services, err := h.Service.GetServices()
	if err != nil {
		log.Println("error getting services:", err)
		return
	}

	m := tgbotapi.NewMessage(msg.Chat.ID, "Выберите услугу 💅")

	// Получаем стандартную клавиатуру услуг
	kb := keyboard.BuildServiceKeyboard(services)

	// КНОПКА НАЗАД В МЕНЮ
	kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📱 В главное меню", "to_main_menu"),
	))

	m.ReplyMarkup = kb

	_, err = h.Bot.Send(m)
	if err != nil {
		log.Println(err)
	}
}
