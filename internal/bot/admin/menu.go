package admin

import (
	"cosmetologybotliza/internal/auth"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
Этот файл - визуальный интерфейс админ-панели.
Его задача — показать пользователю-администратору кнопки управления, если он прошел проверку доступа.
*/
func (h *Handler) OpenAdminMenu(chatID int64, messageID int) {

	// 1. Проверка прав
	user, err := h.User.GetByTelegramID(chatID)
	if err != nil {
		log.Println(err)
		return
	}

	if !auth.IsAdmin(user) {
		_, err = h.Bot.Send(tgbotapi.NewMessage(chatID, "⛔ Нет доступа"))
		if err != nil {
			log.Printf("Ошибка при отправке: %v", err)
		}
		return
	}

	// 2. Текст меню
	text := "🛠 Админ-панель"

	// 3. Формируем клавиатуру
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👩‍💼 Мастера", "admin_masters"),
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои записи", "admin_my_bookings"),
			tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить старые отмены", "admin_clear_cancelled"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ В главное меню", "menu_main"),
		),
	)

	// 4. Отправляем или редактируем
	if messageID > 0 {
		// Редактируем старое сообщение
		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ParseMode = "Markdown"
		edit.ReplyMarkup = &keyboard
		_, err = h.Bot.Send(edit)
		if err != nil {
			log.Printf("Ошибка при отправке: %v", err)
		}
	} else {
		// Шлем новое сообщение
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = keyboard
		_, err = h.Bot.Send(msg)
		if err != nil {
			log.Printf("Ошибка при отправке: %v", err)
		}
	}
}
