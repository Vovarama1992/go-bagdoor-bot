package main

import (
	"log"
	"os"
	"time"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

func main() {
	// Загружаем переменные из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Получаем URL базы данных
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Инициализация базы данных
	pool, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Получаем токен бота из переменных окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	// Настройка бота
	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	// Создаем бота
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("failed to start bot: %v", err)
	}

	// Инициализируем репозиторий с базой данных
	userRepo := user.NewPostgresRepository(pool)

	// Инициализируем сервис, который использует репозиторий
	userService := user.NewService(userRepo)

	// Регистрируем хендлеры, передавая сервис
	bot.Handle("/start", user.HandleStart(userService))

	log.Println("bot running...")

	// Запускаем бота
	bot.Start()
}
