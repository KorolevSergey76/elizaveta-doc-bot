package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"cosmetologybotliza/internal/bot"
	"cosmetologybotliza/internal/bot/admin"
	"cosmetologybotliza/internal/bot/booking"
	"cosmetologybotliza/internal/bot/menu"
	"cosmetologybotliza/internal/config"
	"cosmetologybotliza/internal/fsm"
	"cosmetologybotliza/internal/service"
	"cosmetologybotliza/internal/storage"
)

func main() {
	// 1. Загрузка конфигурации: читаем переменные окружения (токен бота, параметры БД)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}
	cfg := config.NewConfig()

	// 2. Инициализация хранилищ: создание подключений к базам данных
	db := storage.NewPostgres(cfg) // База данных (услуги, пользователи, записи)
	rdb := storage.NewRedis(cfg)   // Кэш/состояние (для FSM, чтобы помнить шаг записи)

	// 3. Инициализация сервисов: бизнес-логика (работа с данными)
	bookingService := &service.BookingService{DB: db} // Логика бронирования
	userService := &service.UserService{DB: db}       // Логика пользователей

	// FSM (Finite State Machine) — для пошаговых диалогов (например: выбор даты -> времени -> услуг)
	fsmService := fsm.NewFSM(rdb)

	// 4. Инициализация Telegram API: устанавливаем связь с сервером Telegram
	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// 5. Инициализация под-хендлеров: создаем объекты, которые будут отвечать за разные части бота
	menuHandler := menu.NewHandlerMenu(botAPI, bookingService, userService)               // Отвечает за меню
	adminHandler := admin.NewHandler(botAPI, bookingService, userService)                 // Отвечает за админку
	bookingHandler := booking.NewHandler(botAPI, bookingService, userService, fsmService) // Отвечает за запись

	// 6. Создание Роутера: "диспетчер", который решает, какой хендлер должен обрабатывать запрос
	router := &bot.Router{
		Menu:    menuHandler,
		Admin:   adminHandler,
		Booking: bookingHandler,
	}

	// 7. Создание главного обработчика: точка входа для всех обновлений из Telegram
	handler := &bot.Handler{
		Bot:    botAPI,
		Router: router,
	}

	// 8. Запуск бота: создаем канал для получения обновлений (сообщений, нажатий кнопок)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	log.Println("Bot started!")
	// Главный цикл: бесконечно слушаем новые сообщения и передаем их в handler
	for update := range updates {
		handler.Handle(update)
	}

	// В main.go, после инициализации бота и сервисов:
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			// Проверяем записи, которые закончились более 1 часа назад
			threshold := time.Now().Add(-1 * time.Hour)
			appointments, err := bookingService.GetFinishedAppointments(threshold)
			if err != nil {
				log.Println("Ошибка при поиске записей для отзыва:", err)
				continue
			}

			for _, app := range appointments {
				// ИСПРАВЛЕНИЕ: вызываем функцию (ее нужно определить ниже)
				booking.SendReviewRequest(botAPI, app)

				// ИСПРАВЛЕНИЕ: используем экземпляр bookingService
				bookingService.MarkReviewSent(app.ID)
			}
		}
	}()
}
