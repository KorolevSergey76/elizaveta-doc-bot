package booking

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/bot/keyboard"
	"cosmetologybotliza/internal/fsm"
)

/*
Этот метод — логический переход от выбора мастера к выбору даты.
Здесь вы связываете выбор пользователя с конечным автоматом (FSM) и подготавливаете интерактивный календарь.
*/

func (h *Handler) selectMaster(
	cb *tgbotapi.CallbackQuery,
	data string,
) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	idStr := strings.TrimPrefix(data, "master_")

	masterID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("master parse error:", err)
		return
	}

	session, err := h.FSM.Get(chatID)
	if err != nil {
		log.Println(err)
		return
	}

	serviceData, err := h.Service.GetServiceByID(session.ServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	session.State = fsm.SelectDate
	session.MasterID = int64(masterID)

	if err := h.FSM.Set(chatID, *session); err != nil {
		log.Println("FSM set error:", err)
		return
	}

	// Генерируем календарь дат
	kb := keyboard.BuildSmartDateKeyboard(
		h.Service,
		int64(masterID),
		serviceData.DurationMinutes,
	)

	// --- Добавляем наши кастомные кнопки управления ---
	kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_to_masters"),
		tgbotapi.NewInlineKeyboardButtonData("📱 В меню", "to_main_menu"),
	))

	replyMarkup := &kb

	edit := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		"Выберите дату 📅",
	)
	edit.ReplyMarkup = replyMarkup

	_, err = h.Bot.Request(edit)
	if err != nil {
		log.Println("edit error:", err)
	}
}
