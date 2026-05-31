package menu

import (
	"cosmetologybotliza/internal/domain"
	"cosmetologybotliza/internal/service"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerMenu struct {
	Bot     *tgbotapi.BotAPI
	Service *service.BookingService
	User    *service.UserService
}

func NewHandlerMenu(bot *tgbotapi.BotAPI, svc *service.BookingService, userSvc *service.UserService) *HandlerMenu {
	return &HandlerMenu{
		Bot:     bot,
		Service: svc,
		User:    userSvc,
	}
}

func (h *HandlerMenu) Handle(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	if strings.HasPrefix(cb.Data, "svc_") {
		// Если это кнопка категории, которую мы еще не обработали,
		// показываем описание из словаря
		if desc, ok := ServiceDescriptions[cb.Data]; ok {
			h.showServiceDetail(chatID, messageID, desc)
			return
		}
	}

	switch cb.Data {

	// Ловит и нажатие из Главного меню, и возврат из визиток мастеров
	case "menu_masters", "menu_back_to_masters":
		h.showMasters(chatID, messageID)

	case "menu_services":
		h.showServices(chatID, messageID)

	case "menu_contacts":
		h.showContacts(chatID, messageID)

	case "menu_location":
		h.showLocation(chatID, messageID)

	case "menu_master_1":
		h.showMaster(chatID, messageID, "master_1")

	case "menu_master_2":
		h.showMaster(chatID, messageID, "master_2")

	case "to_main_menu":
		h.ShowMainMenuEdit(chatID, messageID, cb.From.ID)

	case "menu_education":
		h.showEducation(chatID, messageID)

	case "menu_main":
		h.ShowMainMenuEdit(cb.Message.Chat.ID, cb.Message.MessageID, cb.From.ID)
	}
}

func (h *HandlerMenu) showMasters(chatID int64, messageID int) {
	// 1. Удаляем старое сообщение (с фото мастера)
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, _ = h.Bot.Request(deleteMsg)

	// 2. Создаем новое текстовое сообщение с клавиатурой выбора мастера
	text := "У нас работают топовые специалисты. Выберите мастера, чтобы узнать о нем подробнее: 👩‍💼"
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👑 Королева Елизавета", "menu_master_1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💅 Лапина Полина", "menu_master_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ В главное меню", "to_main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = &kb
	h.Bot.Send(msg)
}

func (h *HandlerMenu) showMaster(chatID int64, messageID int, masterKey string) {
	master, ok := Masters[masterKey]
	if !ok {
		return
	}

	// Удаляем старое
	h.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ К списку мастеров", "menu_back_to_masters"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 В главное меню", "to_main_menu"),
		),
	)

	// Отправляем фото с описанием
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(master.PhotoID))
	photo.Caption = master.Bio
	photo.ParseMode = "Markdown"
	photo.ReplyMarkup = &kb

	h.Bot.Send(photo)
}

func (h *HandlerMenu) showServices(chatID int64, messageID int) {
	text := "💅 Услуги:\n\nВыберите категорию"

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Инъекционная косметология", "svc_inject"),
			tgbotapi.NewInlineKeyboardButtonData("IPL", "svc_ipl"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Чистка лица", "svc_clean"),
			tgbotapi.NewInlineKeyboardButtonData("Уходовые процедуры", "svc_care"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пилинги", "svc_pillings"),
			tgbotapi.NewInlineKeyboardButtonData("Аппаратная косметология", "svc_apparat"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удаление новообразований", "svc_del"),
			tgbotapi.NewInlineKeyboardButtonData("Депилция", "svc_depelation"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💬 Консультация", "svc_consult"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ В главное меню", "to_main_menu"),
		),
	)

	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ReplyMarkup = &kb

	_, _ = h.Bot.Request(msg)
}

func (h *HandlerMenu) ShowMainMenuEdit(chatID int64, messageID int, telegramID int64) {
	// 1. Сначала пробуем удалить старое сообщение (с фото мастера)
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, _ = h.Bot.Request(deleteMsg)

	// 2. Определяем текст и клавиатуру
	text := "💅 Главное меню"
	user, err := h.User.GetByTelegramID(telegramID)
	if err != nil || user == nil {
		user = &domain.User{Role: "client"}
	}
	kb := BuildMainMenu(user)

	// 3. Отправляем НОВОЕ текстовое сообщение
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = &kb

	_, err = h.Bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке главного меню: %v", err)
	}
}

func (h *HandlerMenu) showContacts(chatID int64, messageID int) {
	text := `📞 Контакты:

👩‍💼 Королева Елизавета:
📱 +7 980 664 07 17
✈ dikidi: https://dikidi.net/503969?p=0.pi
📸 Instagram: https://www.instagram.com/koroleva_kosmetolog.yar/pro..
✈ Telegram: @elizavetadoc02
`

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ В главное меню", "to_main_menu"),
		),
	)

	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ReplyMarkup = &kb

	_, _ = h.Bot.Request(msg)
}

func (h *HandlerMenu) showLocation(chatID int64, messageID int) {
	text := `📍 Как добраться:

Адрес: г.Ярославль, ул. Первомайская, 53 (второй этаж)

🚇 Ближайшая остановка: Богоявленская площадь

`

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ В главное меню", "to_main_menu"),
		),
	)

	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ReplyMarkup = &kb

	_, _ = h.Bot.Request(msg)
}

// 1. Меню обучения
func (h *HandlerMenu) showEducation(chatID int64, messageID int) {
	text := "🎓 *Обучение и повышение квалификации*\n\n" +
		"Здесь вы можете пройти профессиональные курсы по косметологии."

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➡️ Продолжить", "edu_blocks"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ В меню", "to_main_menu"),
		),
	)

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ParseMode = "Markdown"
	edit.ReplyMarkup = &kb
	h.Bot.Request(edit)
}

func (h *HandlerMenu) showServiceDetail(chatID int64, messageID int, description string) {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к услугам", "menu_services"),
		),
	)

	msg := tgbotapi.NewEditMessageText(chatID, messageID, description)
	msg.ReplyMarkup = &kb
	h.Bot.Request(msg)
}
