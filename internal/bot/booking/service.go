package booking

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/bot/keyboard"
	"cosmetologybotliza/internal/fsm"
)

// selectService — теперь отображает карточку услуги, а не переходит сразу к мастерам
func (h *Handler) selectService(cb *tgbotapi.CallbackQuery, data string) {
	serviceIDStr := strings.TrimPrefix(data, "service_")
	serviceID, _ := strconv.ParseInt(serviceIDStr, 10, 64)

	h.ShowServiceDetails(cb.Message.Chat.ID, cb.Message.MessageID, serviceID)
}

// ShowServiceDetails — показывает детальную информацию об услуге
func (h *Handler) ShowServiceDetails(chatID int64, messageID int, serviceID int64) {
	service, err := h.Service.GetServiceByID(serviceID)
	if err != nil {
		log.Println("Ошибка получения услуги:", err)
		return
	}

	// Обработка отсутствия описания
	desc := "Описание уточняйте у администратора."
	if service.Description != nil && *service.Description != "" {
		desc = *service.Description
	}

	text := fmt.Sprintf("💆‍♀️ **%s**\n\n📝 %s\n\n⏱ Длительность: %d мин.\n💰 Цена: %d руб.",
		service.Name, desc, service.DurationMinutes, service.Price)

	// Кнопка записи переводит на выбор мастера
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Записаться", fmt.Sprintf("book_confirm_%d", service.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_services"),
		),
	)

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ReplyMarkup = &keyboard
	edit.ParseMode = "Markdown"

	_, err = h.Bot.Request(edit)
	if err != nil {
		log.Println("Ошибка редактирования сообщения:", err)
	}
}

// HandleBookingConfirm — вызывается при нажатии "Записаться" в карточке услуги
func (h *Handler) HandleBookingConfirm(cb *tgbotapi.CallbackQuery) {
	serviceIDStr := strings.TrimPrefix(cb.Data, "book_confirm_")
	serviceID, _ := strconv.ParseInt(serviceIDStr, 10, 64)

	session, err := h.FSM.Get(cb.Message.Chat.ID)
	if err != nil {
		return
	}

	session.State = fsm.SelectMaster
	session.ServiceID = serviceID
	_ = h.FSM.Set(cb.Message.Chat.ID, *session)

	// Показываем список мастеров
	edit := tgbotapi.NewEditMessageText(cb.Message.Chat.ID, cb.Message.MessageID, "Выберите мастера 👩‍💼")
	kb := keyboard.BuildMasterKeyboard(h.Service)
	kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_services"),
	))

	edit.ReplyMarkup = &kb
	h.Bot.Request(edit)
}
