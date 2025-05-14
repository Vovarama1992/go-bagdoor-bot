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
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	pool, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("failed to start bot: %v", err)
	}

	userRepo := user.NewPostgresRepository(pool)

	// маршруты бота
	bot.Handle("/start", user.HandleStart(userRepo))

	log.Println("bot running...")
	bot.Start()
}
