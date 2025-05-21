package user

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

func HandleStart(s *Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		from := c.Sender()
		user, err := s.GetByTgID(ctx, from.ID)
		if err == nil && user != nil {
			return c.Send("Добро пожаловать! 🔥\nРады видеть вас снова", openAppMarkup())
		}

		return c.Send(`Привет! Это бот Bagdoor⚡

Перед началом проверь, что ты подписался(а) на наш [чат](https://t.me/+s4aQ9RU-K9JkZmNi) и [канал](https://t.me/bagdoor) — это необходимо для продолжения.

Продолжая, ты автоматически соглашаешься с Пользовательским соглашением и Политикой конфиденциальности.`, subscribeMarkup())
	}
}

// Хендлер для проверки подписки и запроса номера телефона
func SubscribeHandler(s *Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		// Проверка, подписан ли пользователь на чат и канал
		if !isSubscribed(c) {
			return c.Send("Подписки не найдены! Пожалуйста, подпишитесь на чат и канал.")
		}

		// Отправляем кнопку для получения номера телефона
		return c.Send(
			"Теперь, чтобы завершить регистрацию, отправь свой номер телефона.",
			&tele.SendOptions{ReplyMarkup: phoneMarkup()},
		)
	}
}
func PhoneHandler(s *Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		contact := c.Message().Contact

		if contact.UserID != c.Sender().ID {
			return c.Send("Пожалуйста, отправьте свой собственный номер.")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		newUser := &User{
			TgID:        c.Sender().ID,
			TgUsername:  c.Sender().Username,
			FirstName:   c.Sender().FirstName,
			LastName:    c.Sender().LastName,
			PhoneNumber: contact.PhoneNumber,
		}

		if err := s.RegisterUser(ctx, newUser); err != nil {
			log.Printf("Ошибка при регистрации: %v", err)
			return c.Send("Ошибка при регистрации.")
		}

		// Убираем клавиатуру
		_ = c.Send("✅", &tele.SendOptions{ReplyMarkup: &tele.ReplyMarkup{}})

		// Приветственное сообщение
		return c.Send("Добро пожаловать! 🔥\nТеперь вы можете пользоваться ботом Bagdoor: размещать рейсы, заказы и совершать безопасные сделки. Удачи!",
			&tele.SendOptions{ReplyMarkup: &tele.ReplyMarkup{}},
		)
	}
}

func isSubscribed(c tele.Context) bool {
	bot := c.Bot()
	user := c.Sender()

	chatID := os.Getenv("BAGDOOR_CHAT_ID")       // группа
	channelID := os.Getenv("BAGDOOR_CHANNEL_ID") // канал

	if chatID == "" || channelID == "" {
		log.Println("Не заданы BAGDOOR_CHAT_ID или BAGDOOR_CHANNEL_ID в .env")
		return false
	}

	// Проверка канала
	channelInt, err := strconv.ParseInt(channelID, 10, 64)
	if err != nil {
		log.Printf("Ошибка парсинга BAGDOOR_CHANNEL_ID: %v", err)
		return false
	}
	channel := &tele.Chat{ID: channelInt}
	member, err := bot.ChatMemberOf(channel, user)
	if err != nil {
		log.Printf("Ошибка проверки канала: %v", err)
		return false
	}
	if member.Role == "left" {
		log.Println("Пользователь не подписан на канал")
		return false
	}

	// Проверка чата
	chatInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Printf("Ошибка парсинга BAGDOOR_CHAT_ID: %v", err)
		return false
	}
	chat := &tele.Chat{ID: chatInt}
	member, err = bot.ChatMemberOf(chat, user)
	if err != nil {
		log.Printf("Ошибка проверки чата: %v", err)
		return false
	}
	if member.Role == "left" {
		log.Println("Пользователь не подписан на чат")
		return false
	}

	return true
}

// Функция для создания разметки кнопок для подписки
func subscribeMarkup() *tele.ReplyMarkup {
	btnSubscribed := tele.Btn{Text: "✅ Я подписался"}
	markup := &tele.ReplyMarkup{}
	markup.Reply(markup.Row(btnSubscribed))
	return markup
}

// Функция для создания кнопки с ссылкой на мини-приложение
func openAppMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}
	btn := markup.URL("Открыть", "https://t.me/bagdoorapp_bot/app")
	markup.Reply(markup.Row(btn))
	return markup
}

// Функция для кнопки отправки номера телефона
func phoneMarkup() *tele.ReplyMarkup {
	btnPhone := tele.Btn{Contact: true, Text: "📱 Отправить номер"}

	markup := &tele.ReplyMarkup{}
	markup.Reply(
		markup.Row(btnPhone),
	)
	return markup
}
