package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Bot    *tgbotapi.BotAPI
	Router *Router
}

func (h *Handler) Handle(update tgbotapi.Update) {
	if update.Message != nil {
		h.Router.RouteMessage(update.Message)
		return
	}

	if update.CallbackQuery != nil {
		// Достаточно одного вызова роутера
		h.Router.RouteCallback(update.CallbackQuery)
	}
}
