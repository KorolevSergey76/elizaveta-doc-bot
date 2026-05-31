package booking

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) cancelBooking(
	chatID int64,
	data string,
	cb *tgbotapi.CallbackQuery,
) {

	_, err := h.Bot.Request(tgbotapi.NewCallback(cb.ID, ""))
	if err != nil {
		// Логируем ошибку только если это действительно что-то важное, а не просто "устаревший запрос"
		if !strings.Contains(err.Error(), "query is too old") && !strings.Contains(err.Error(), "not found") {
			log.Printf("Ошибка ответа на callback: %v", err)
		}
	}

	idStr := strings.TrimPrefix(data, "cancel_")
	appointmentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	// 2. Получаем запись перед отменой, чтобы знать детали (кто мастер, на когда)
	appointment, err := h.Service.GetAppointmentByID(appointmentID)
	if err != nil {
		log.Println("Не удалось найти запись для отмены:", err)
		return
	}

	// 3. Выполняем отмену
	user, err := h.UserService.GetByTelegramID(chatID)
	if err != nil || user == nil {
		log.Println("Ошибка получения пользователя:", err)
		return
	}
	err = h.Service.CancelAppointment(user.ID, appointmentID)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка отмены ❌"))
		return
	}

	// 4. Уведомляем мастера
	// Предположим, у вас есть метод получения TelegramID мастера по его ID
	masterTelegramID := h.UserService.GetTelegramIDByMasterID(appointment.MasterID)
	if masterTelegramID != 0 {
		msgText := fmt.Sprintf("⚠️ **Запись отменена клиентом!**\n\n👤 Клиент: %s\n🕒 Время: %s",
			user.Username, appointment.Time.Format("02.01.2006 15:04"))

		msg := tgbotapi.NewMessage(masterTelegramID, msgText)
		msg.ParseMode = "Markdown"
		h.Bot.Send(msg)
	}

	// 5. Обновляем сообщение клиенту
	edit := tgbotapi.NewEditMessageText(chatID, cb.Message.MessageID, "Ваша запись была успешно отменена ✅")
	h.Bot.Send(edit)
	h.Bot.Request(tgbotapi.NewCallback(cb.ID, "Отменено"))
}
