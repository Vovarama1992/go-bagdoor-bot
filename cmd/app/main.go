package main

import (
	"log"
	"os"
	"time"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

func main() {
	// Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Инициализация базы данных
	pool, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
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
		log.Fatalf("failed to start bot: %v", err)
	}

	userRepo := user.NewPostgresRepository(pool)
	userService := user.NewService(userRepo)

	bot.Handle("/start", user.HandleStart(userService))
	bot.Handle("✅ Я подписался", user.SubscribeHandler(userService))
	bot.Handle(tele.OnContact, user.PhoneHandler(userService))
        bot.Handle("/getchannelid", func(c tele.Context) error {
	chat, err := bot.ChatByUsername("@bagdoor") // или другой username
	if err != nil {
		log.Printf("Ошибка получения канала: %v", err)
		return c.Send("Ошибка получения канала.")
	}
	log.Printf("Channel ID: %d", chat.ID)
	return c.Send("Channel ID зафиксирован.")
})

	log.Println("bot running...")
	bot.Start()
}
