package main

import (
	"log"
	"os"
	"time"

	botflight "github.com/Vovarama1992/go-bagdoor-bot/internal/bot/flight"
	botorder "github.com/Vovarama1992/go-bagdoor-bot/internal/bot/order"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/flight"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/storage"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

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

	userRepo := user.NewPostgresRepository(pool)
	userService := user.NewService(userRepo)

	orderRepo := order.NewPostgresRepository(pool)
	orderService := order.NewService(orderRepo, userRepo)

	flightRepo := flight.NewPostgresRepository(pool)
	flightService := flight.NewService(flightRepo)

	s3Uploader := storage.NewS3Uploader()

	bot.Handle("/start", user.HandleStart(userService))
	bot.Handle(tele.OnText, func(c tele.Context) error {
		if c.Text() == "✅ Я подписался" {
			return user.SubscribeHandler(userService)(c)
		}
		return nil
	})
	bot.Handle(tele.OnContact, user.PhoneHandler(userService))

	bot.Handle("/setorderid", botorder.HandleSetOrderID())
	bot.Handle(tele.OnPhoto, botorder.HandlePhotoUpload(orderService, bot, s3Uploader))

	bot.Handle("/setflightid", botflight.HandleSetFlightID())
	bot.Handle(tele.OnDocument, botflight.HandlePdfUpload(flightService, bot, s3Uploader))

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
