package bot

import (
	"cosmetologybotliza/internal/bot/admin"
	"cosmetologybotliza/internal/bot/booking"
	"cosmetologybotliza/internal/bot/menu"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	Menu    *menu.HandlerMenu
	Booking *booking.Handler
	Admin   *admin.Handler
}

func (r *Router) RouteCallback(cb *tgbotapi.CallbackQuery) {
	data := cb.Data

	if strings.HasPrefix(data, "menu_") || data == "to_main_menu" {
		r.Menu.Handle(cb)
		return
	}

	if strings.HasPrefix(data, "menu_") || strings.HasPrefix(data, "svc_") || data == "to_main_menu" {
		r.Menu.Handle(cb)
		return
	}

	if strings.HasPrefix(data, "admin_") {
		r.Admin.HandleCallback(cb)
		return
	}

	r.Booking.HandleCallback(cb)
}

func (r *Router) RouteMessage(msg *tgbotapi.Message) {
	// 1. Сначала проверяем, является ли сообщение командой
	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			// Вызываем ваш существующий метод
			// 0 вместо messageID, так как при старте у нас нет ID старого сообщения
			r.Menu.ShowMainMenuEdit(msg.Chat.ID, 0, msg.From.ID)
			return
		case "admin":
			r.Admin.HandleMessage(msg)
			return
		}
	}

	// 2. Если это не команда, обрабатываем как обычное текстовое сообщение
	// (например, для ввода данных бронирования)
	r.Booking.HandleMessage(msg)
}
