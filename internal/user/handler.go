package user

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

var BtnConfirmSub = &tele.Btn{Unique: "confirm_sub"}

func HandleStart(s *Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		from := c.Sender()
		user, err := s.GetByTgID(ctx, from.ID)
		if err == nil && user != nil {
			log.Printf("👋 Повторный вход: пользователь %d уже зарегистрирован", from.ID)
			return c.Send("Добро пожаловать! 🔥\nРады видеть вас снова", openAppMarkup())
		}

		_ = c.Send("...", &tele.SendOptions{
			ReplyMarkup: &tele.ReplyMarkup{RemoveKeyboard: true},
		})

		// Инлайн-кнопки
		markup := &tele.ReplyMarkup{}
		btnChat := markup.URL("🔗 Чат", "https://t.me/+s4aQ9RU-K9JkZmNi")
		btnChannel := markup.URL("📣 Канал", "https://t.me/bagdoor")
		btnConfirm := markup.Data("✅ Я подписался", BtnConfirmSub.Unique)

		markup.Inline(
			markup.Row(btnChat, btnChannel),
			markup.Row(btnConfirm),
		)

		return c.Send("Привет! Это бот Bagdoor⚡\n\nПеред началом подпишись на чат и канал, а затем нажми '✅ Я подписался'", markup)
	}
}

func SubscribeHandler(s *Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		log.Printf("Нажата кнопка '✅ Я подписался' от пользователя %d", c.Sender().ID)
		_ = c.Respond()

		if !isSubscribed(c) {
			log.Printf("❌ Пользователь %d не подписан", c.Sender().ID)
			return c.Send("Подписки не найдены! Пожалуйста, подпишитесь на чат и канал.")
		}

		log.Printf("✅ Пользователь %d успешно прошёл проверку подписки", c.Sender().ID)

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
			log.Printf("❌ Ошибка при регистрации: %v", err)
			return c.Send("Ошибка при регистрации.")
		}

		log.Printf("✅ Пользователь %d успешно зарегистрирован", c.Sender().ID)

		opts := &tele.SendOptions{ReplyMarkup: &tele.ReplyMarkup{RemoveKeyboard: true}}

		_ = c.Send("✅", opts)

		log.Printf("➡️ Отправляем приветствие и кнопку на миниаппу пользователю %d", c.Sender().ID)

		return c.Send(
			"Добро пожаловать! 🔥\nТеперь вы можете пользоваться ботом Bagdoor: размещать рейсы, заказы и совершать безопасные сделки. Удачи!",
			openAppMarkup(),
		)
	}
}

func isSubscribed(c tele.Context) bool {
	bot := c.Bot()
	user := c.Sender()

	chatID := os.Getenv("BAGDOOR_CHAT_ID")
	channelID := os.Getenv("BAGDOOR_CHANNEL_ID")

	if chatID == "" || channelID == "" {
		log.Println("❗ Не заданы BAGDOOR_CHAT_ID или BAGDOOR_CHANNEL_ID в .env")
		return false
	}

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

func openAppMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}
	btn := tele.Btn{
		Text: "Открыть",
		WebApp: &tele.WebApp{
			URL: "https://tgbot.bagdoor.io",
		},
	}
	markup.Inline(markup.Row(btn))
	return markup
}

func phoneMarkup() *tele.ReplyMarkup {
	btnPhone := tele.Btn{Contact: true, Text: "📱 Отправить номер"}

	markup := &tele.ReplyMarkup{}
	markup.Reply(
		markup.Row(btnPhone),
	)
	return markup
}
