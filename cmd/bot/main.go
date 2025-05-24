package main

import (
	"log"
	"os"
	"time"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/flight"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/storage"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(dbURL string) {
	m, err := migrate.New("file://internal/db/migrations", dbURL)
	if err != nil {
		log.Fatalf("Ошибка инициализации миграций: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка при применении миграций: %v", err)
	}
	log.Println("Миграции успешно применены")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	runMigrations(dbURL)

	pool, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе: %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	pref := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}

	// Инициализация сервисов
	userRepo := user.NewPostgresRepository(pool)
	userService := user.NewService(userRepo)

	orderRepo := order.NewPostgresRepository(pool)
	flightRepo := flight.NewPostgresRepository(pool)
	orderService := order.NewService(orderRepo, userRepo)
	flightService := flight.NewService(flightRepo)

	s3Uploader := storage.NewS3Uploader()

	// Хендлеры пользователя
	bot.Handle("/start", user.HandleStart(userService))
	bot.Handle("✅ Я подписался", user.SubscribeHandler(userService))
	bot.Handle(tele.OnContact, user.PhoneHandler(userService))

	// Хендлеры заказов
	bot.Handle("/setorderid", order.HandleSetOrderID())
	bot.Handle(tele.OnPhoto, order.HandlePhotoUpload(orderService, bot, s3Uploader))
	bot.Handle("/setflightid", flight.HandleSetFlightID())
	bot.Handle(tele.OnDocument, flight.HandlePdfUpload(flightService, bot, s3Uploader))

	// Отладка: получить ID канала
	bot.Handle("/getchannelid", func(c tele.Context) error {
		chat, err := bot.ChatByUsername("@bagdoor")
		if err != nil {
			log.Printf("Ошибка получения канала: %v", err)
			return c.Send("Ошибка получения канала.")
		}
		log.Printf("Channel ID: %d", chat.ID)
		return c.Send("Channel ID зафиксирован.")
	})

	log.Println("Bot running...")
	bot.Start()
}
