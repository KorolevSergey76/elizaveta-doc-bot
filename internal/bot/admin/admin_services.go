package admin

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleClearCancelled(cb *tgbotapi.CallbackQuery) {
	count, err := h.Service.DeleteOldCancelledAppointments(30) // Удаляем всё старше 30 дней
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "Ошибка удаления ❌"))
		return
	}
	h.Bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, fmt.Sprintf("Успешно удалено записей: %d ✅", count)))
}
