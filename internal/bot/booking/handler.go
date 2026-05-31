package booking

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/fsm"
	"cosmetologybotliza/internal/service"
)

func NewHandler(
	bot *tgbotapi.BotAPI,
	svc service.BookingServiceInterface,
	userSvc *service.UserService,
	fsmSvc *fsm.FSM,
) *Handler {
	return &Handler{
		Bot:         bot,
		Service:     svc,
		UserService: userSvc,
		FSM:         fsmSvc,
	}
}

// entry point для bot handler
func (h *Handler) HandleCallback(cb *tgbotapi.CallbackQuery) {
	h.routeCallback(cb)
}

func (h *Handler) HandleMessage(msg *tgbotapi.Message) {

	switch msg.Text {
	case "/start", "📱 Главное меню":
		h.Start(msg)
	case "/book":
		h.StartBooking(msg)
	case "/cancel":
		_ = h.FSM.Clear(msg.Chat.ID)
		h.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Запись отменена ❌"))
	case "/mybookings":
		h.ShowUserBookings(msg.Chat.ID)
	}

}
