package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"cosmetologybotliza/internal/bot"
	"cosmetologybotliza/internal/repository"
	"cosmetologybotliza/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Драйвер базы данных
)

func main() {
	// 1. Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные среды системы")
	}

	// 2. Инициализация БД
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panic("Ошибка при попытке открыть БД: ", err)
	}
	defer db.Close()

	// Проверка связи с БД
	if err := db.Ping(); err != nil {
		log.Panic("База данных недоступна: ", err)
	}
	log.Println("Подключение к БД успешно установлено!")

	// 3. Инициализация API бота
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан в .env")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Ошибка инициализации бота: ", err)
	}

	// 4. Сборка слоев (Dependency Injection)
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	handler := &bot.Handler{
		Bot:     api,
		Service: svc,
	}

	// 5. Запуск получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := api.GetUpdatesChan(u)

	log.Printf("Бот успешно запущен как @%s", api.Self.UserName)

	for update := range updates {
		// Передаем обновление в наш обработчик
		handler.HandleUpdate(update)
	}
}
