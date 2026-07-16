package bot

import (
	"cosmetologybotliza/internal/bot/menu"
	"cosmetologybotliza/internal/service"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Bot     *tgbotapi.BotAPI
	Service *service.Service
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {

	/*if update.Message != nil && len(update.Message.Photo) > 0 {
		photos := update.Message.Photo
		fileID := photos[len(photos)-1].FileID
		log.Printf("🔥 СВЕЖИЙ FILE_ID: %s", fileID)
	}*/

	if update.Message != nil {
		h.handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		h.handleCallback(update.CallbackQuery)
	}
}

func (h *Handler) handleMessage(msg *tgbotapi.Message) {
	if msg.IsCommand() && msg.Command() == "start" {
		h.Service.SaveUser(msg.From.ID, msg.From.UserName, msg.From.FirstName)

		reply := tgbotapi.NewMessage(msg.Chat.ID, menu.WelcomeText)
		reply.ReplyMarkup = menu.MainMenuKeyboard()
		h.Bot.Send(reply)
	}
}

func (h *Handler) handleCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	// Логируем, что пришло
	log.Printf("Callback data: %s", cb.Data)

	switch cb.Data {
	case "menu_main":
		h.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
		h.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID+1)) // Удаляем карту, если она была отправлена следом

		newMsg := tgbotapi.NewMessage(chatID, menu.WelcomeText)
		newMsg.ReplyMarkup = menu.MainMenuKeyboard()
		h.Bot.Send(newMsg)

	case "menu_services":
		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите услугу для подробного ознакомления:")
		kb := menu.ServicesCategoriesKeyboard()
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "service_lips_down", "service_botox", "service_contour",
		"service_bio", "service_hardware", "service_consultation", "service_care",
		"service_peeling", "service_cleaning", "service_depilation":

		text := menu.ServicesDescriptions[cb.Data]
		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ParseMode = "Markdown"
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("🗓 Записаться", "https://dikidi.net/ru/profile/kosmetolog_korolyova_elizaveta_503969")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_services")),
		)
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "menu_doctor":
		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Выберите врача:")
		kb := menu.DoctorListKeyboard()
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "doctor_eliza", "doctor_polina":
		h.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
		h.sendDoctorProfile(chatID, cb.Data)

	case "back_to_doctors":
		// 1. Удаляем текущее сообщение (текст с кнопкой "Назад")
		h.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))

		// 2. Удаляем предыдущее сообщение (фото).
		// Обычно фото — это messageID - 1
		h.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID-1))

		// 3. Отправляем новое меню врачей
		newMsg := tgbotapi.NewMessage(chatID, "Выберите врача:")
		newMsg.ReplyMarkup = menu.DoctorListKeyboard()
		h.Bot.Send(newMsg)

	case "menu_location":
		h.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
		// Отправляем текст
		msg := tgbotapi.NewMessage(chatID, menu.LocationText)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = menu.BackButton()
		h.Bot.Send(msg)
		// Отправляем карту
		h.Bot.Send(tgbotapi.NewLocation(chatID, 57.623024, 39.889422))

	case "menu_contacts":
		h.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))

		msg := tgbotapi.NewMessage(chatID, "Выберите способ связи:")
		msg.ReplyMarkup = menu.ContactsKeyboard()
		h.Bot.Send(msg)

	case "show_phone":
		// Это действие просто показывает номер через callback-ответ
		callbackCfg := tgbotapi.NewCallback(cb.ID, "Телефон: +7 (980) 664-07-17")
		callbackCfg.ShowAlert = true // Появится всплывающее окно
		h.Bot.Request(callbackCfg)

	case "service_lips_up":
		text := "Как ни крути, но существуют противопоказания к процедуре увеличения губ - пройди короткий тест прямо сейчас, чтобы избавиться от всех сомнений ⤵️"
		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пройти тест", "quiz_q1")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("👄 Памятка при увеличении губ", "menu_lips_memo")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_services")),
		)
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "menu_lips_memo":
		edit := tgbotapi.NewEditMessageText(chatID, messageID, menu.LipsMemoText)
		edit.ParseMode = "Markdown"
		// Возвращаем в главное меню
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к услуге", "service_lips_up")),
		)
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "menu_reviews":
		edit := tgbotapi.NewEditMessageText(chatID, messageID, menu.ReviewsText)
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_main")))
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "menu_education":
		edit := tgbotapi.NewEditMessageText(chatID, messageID, menu.EducationText)
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_main")))
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	// Обработка квиза
	case "quiz_q1":
		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Вопрос 1: Это будет первичная коррекция губ в вашей жизни?")
		kb := menu.LipsQuizKeyboard("q1")
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "quiz_q1_yes", "quiz_q1_no":
		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Вопрос 2: Есть ли у вас аллергия на лидокаин?")
		kb := menu.LipsQuizKeyboard("q2")
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "quiz_q2": // Переход ко второму вопросу
		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Вопрос 2: Есть ли аллергия на лидокаин?")
		kb := menu.LipsQuizKeyboard("q2")
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "quiz_q2_no": // Нет аллергии
		h.sendBookingSuccess(chatID, messageID)

	case "quiz_q2_yes": // Есть аллергия
		h.sendContactDoctor(chatID, messageID)

	case "quiz_q2_idk": // Не знаю про аллергию
		edit := tgbotapi.NewEditMessageText(chatID, messageID, "Лечили ли вы когда-нибудь зубы у стоматолога с анестезией?")
		kb := menu.LipsQuizKeyboard("q3")
		edit.ReplyMarkup = &kb
		h.Bot.Send(edit)

	case "quiz_q3_yes": // Лечили зубы (все ок)
		h.sendBookingSuccess(chatID, messageID)

	case "quiz_q3_no", "quiz_q3_idk": // Не лечили/не помнят
		h.sendContactDoctor(chatID, messageID)
	}

	h.Bot.Request(tgbotapi.NewCallback(cb.ID, ""))
}

// sendDoctorProfile отправляет фото и описание врача
func (h *Handler) sendDoctorProfile(chatID int64, doctorID string) {
	var fileID, caption string

	if doctorID == "doctor_eliza" {
		fileID = "AgACAgIAAxkDAANrah0hPlnrg7MBgIaAsqy3KPePN_UAAx1rG3kQ4Ehd_2DmEYZFCAEAAwIAA3kAAzsE"
		caption = "👩‍⚕️ **Елизавета Королева — врач-косметолог**\n\n" +
			"Я — врач-косметолог, который помогает пациентам сохранять молодость, здоровье кожи и уверенность в своей внешности.\n\n" +
			"В своей работе я сочетаю медицинский подход, современные технологии и внимательное отношение. Моя цель — подчеркнуть вашу естественную красоту, сохранив индивидуальность и гармонию.\n\n" +
			"🎓 **Образование и опыт:**\n" +
			"Ярославский государственный медицинский университет (Лечебное дело). Регулярно повышаю квалификацию, внедряя в практику только проверенные и безопасные протоколы.\n\n" +
			"💉 **Спектр услуг:**\n" +
			"• Контурная пластика и ботулинотерапия\n" +
			"• Биоревитализация, мезо- и плазмотерапия\n" +
			"• Коллагеностимуляция и аппаратная косметология\n" +
			"• Лечение акне, розацеа, пигментации, дерматитов\n" +
			"• Коррекция возрастных изменений лица и тела\n" +
			"• Индивидуальные программы омоложения\n\n" +
			"🎓 **Обучение:**\n" +
			"Я обучаю врачей-косметологов, делясь практическим опытом и авторскими техниками.\n\n" +
			"✨ **Почему мне доверяют:**\n" +
			"Профессионализм, честный подход (ничего лишнего!) и глубокая проработка индивидуальных программ с учетом вашего образа жизни.\n\n" +
			"Приходите на консультацию — составим ваш план красоты!"
	} else {
		fileID = "AgACAgIAAxkBAAO8ah2VdFm6gKMAAajhERVjwU5PqG1iAAKlHGsbeRDoSKpXWhZcRCr6AQADAgADeQADOwQ"
		caption = "🌸 **Полина Лапина — Косметолог-эстетист**\n\n" +
			"Я — косметолог-эстетист с высшим медицинским образованием. Помогаю сохранить естественную красоту кожи, подчеркнуть ухоженность и вернуть лицу здоровое сияние без агрессивных вмешательств.\n\n" +
			"В работе делаю акцент на мягком, деликатном и системном подходе. Моя цель — улучшить качество кожи, восстановить её баланс и подчеркнуть вашу природную привлекательность.\n\n" +
			"🎓 **Опыт и знания:**\n" +
			"Постоянно совершенствую навыки в эстетической косметологии, осваивая самые современные и безопасные методики для лица и тела.\n\n" +
			"✨ **Мои услуги:**\n" +
			"• Профессиональные чистки лица\n" +
			"• Пилинги различной глубины воздействия\n" +
			"• Уходовые и восстановительные протоколы\n" +
			"• Массажи лица (лимфодренажные, лифтинг-техники)\n" +
			"• Уход за проблемной кожей (акне, постакне)\n" +
			"• Подбор индивидуального домашнего ухода\n" +
			"• Профилактические anti-age процедуры\n\n" +
			"🤍 **Мой подход:**\n" +
			"Внимательность, аккуратность и бережность. Для меня важно не просто выполнить процедуру, а выстроить систему ухода, которая дает устойчивый и заметный результат.\n\n" +
			"Помогу вашей коже «зазвучать заново» через грамотный и профессиональный уход!"
	}

	// 1. Отправляем фото
	photoMsg, err := h.Bot.Send(tgbotapi.NewPhoto(chatID, tgbotapi.FileID(fileID)))
	if err != nil {
		log.Printf("Ошибка отправки фото: %v", err)
		return
	}
	log.Printf("Фото отправлено, ID: %d", photoMsg.MessageID)

	// 2. Отправляем описание
	msg := tgbotapi.NewMessage(chatID, caption)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к врачам", "back_to_doctors"),
		),
	)

	textMsg, err := h.Bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки текста: %v", err)
	} else {
		log.Printf("Текст отправлен, ID: %d", textMsg.MessageID)
	}
}

func (h *Handler) sendBookingSuccess(chatID int64, messageID int) {
	text := "Отлично! Противопоказаний не выявлено. Вы можете записаться на процедуру онлайн:"
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("🗓 Записаться", "https://dikidi.net/ru/profile/kosmetolog_korolyova_elizaveta_503969")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_services")),
	)
	edit.ReplyMarkup = &kb
	h.Bot.Send(edit)
}

func (h *Handler) sendContactDoctor(chatID int64, messageID int) {
	text := "Для безопасности нам нужно обсудить это индивидуально. Напишите мне в личные сообщения!"
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("💬 Написать врачу", "https://t.me/@elizavetadoc02")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu_services")),
	)
	edit.ReplyMarkup = &kb
	h.Bot.Send(edit)
}
