package booking

import (
	"cosmetologybotliza/internal/domain"
	"cosmetologybotliza/internal/service"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Планировщик, который проверяет записи
func StartReviewWorker(bot *tgbotapi.BotAPI, service *service.BookingService) {
	ticker := time.NewTicker(1 * time.Hour) // Проверяем раз в час
	for range ticker.C {
		// Ищем записи, которые закончились 1-2 часа назад
		appointments, err := service.GetFinishedAppointments(time.Now().Add(-2 * time.Hour))
		if err != nil {
			continue
		}

		for _, app := range appointments {
			SendReviewRequest(bot, app)
			// Помечаем в БД, что запрос на отзыв отправлен, чтобы не спамить
			service.MarkReviewSent(app.ID)
		}
	}
}

func SendReviewRequest(bot *tgbotapi.BotAPI, app domain.Appointment) {
	text := fmt.Sprintf("💅 *%s*, как прошла ваша процедура у мастера %s?\n\nОцените результат:",
		app.ClientName, app.MasterName)

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", "rate_1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "rate_2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "rate_3"),
			tgbotapi.NewInlineKeyboardButtonData("4", "rate_4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "rate_5"),
		),
	)

	msg := tgbotapi.NewMessage(app.ClientTelegramID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = kb
	bot.Send(msg)
}

func (h *Handler) handleRate(cb *tgbotapi.CallbackQuery) {
	rating := strings.TrimPrefix(cb.Data, "rate_")

	if rating == "5" || rating == "4" {
		// Если оценка высокая — просим оставить отзыв на картах
		text := "Спасибо! Нам очень важно ваше мнение. Будем рады отзыву на Яндекс.Картах: [ссылка]"
		h.Bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, text))
	} else {
		// Если низкая — просим написать, что не понравилось
		text := "Нам очень жаль, что вы остались недовольны. Пожалуйста, напишите, что мы можем улучшить?"
		h.Bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, text))

		// Тут можно перевести FSM в состояние "WaitingForFeedback"
		// чтобы следующее сообщение пользователя пришло к вам как жалоба
	}
}
