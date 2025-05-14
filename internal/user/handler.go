package user

import (
	tele "gopkg.in/telebot.v3"
)

// Стартовый хендлер для проверки tgID и предложения подписки
func StartHandler(s *Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		from := c.Sender()

		// Проверка, существует ли пользователь с данным tgID
		ctx := c
		user, err := s.GetByTgID(ctx, from.ID)
		if err == nil && user != nil {
			// Пользователь найден, приветственное сообщение
			return c.Send(
				"Добро пожаловать! 🔥\nРады видеть вас снова",
				openAppMarkup(),
			)
		}

		// Новый пользователь, предложим подписаться
		return c.Send(
			`Привет! Это бот Bagdoor⚡

Перед началом проверь, что ты подписан(а) на наш [чат](https://t.me/+s4aQ9RU-K9JkZmNi) и [канал](https://t.me/bagdoor) — это необходимо для продолжения.

Продолжая, ты автоматически соглашаешься с Пользовательским соглашением и Политикой конфиденциальности.`,
			subscribeMarkup(),
		)
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

// Функция для обработки отправки номера телефона
func PhoneHandler(s *Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		contact := c.Message().Contact

		// Проверяем, что отправленный контакт принадлежит пользователю
		if contact.UserID != c.Sender().ID {
			return c.Send("Пожалуйста, отправьте свой собственный номер.")
		}

		// Получаем контекст
		ctx := c

		// Создаем нового пользователя с данными
		newUser := &User{
			TgID:        c.Sender().ID,
			TgUsername:  c.Sender().Username,
			FirstName:   c.Sender().FirstName,
			LastName:    c.Sender().LastName,
			PhoneNumber: contact.PhoneNumber,
		}

		// Регистрируем нового пользователя
		if err := s.RegisterUser(ctx, newUser); err != nil {
			return c.Send("Ошибка при регистрации.")
		}

		// Приветственное сообщение после завершения регистрации
		return c.Send(
			"Добро пожаловать! 🔥\nТеперь вы можете пользоваться ботом Bagdoor: размещать рейсы, заказы и совершать безопасные сделки. Удачи!",
			openAppMarkup(),
		)
	}
}

// Функция для проверки подписки
func isSubscribed(c tele.Context) bool {
	// Здесь должна быть реальная логика проверки подписки через Telegram API
	// Временно возвращаем true, чтобы пользователи могли продолжить регистрацию.
	return true
}

// Функция для создания разметки кнопок для подписки
func subscribeMarkup() *tele.ReplyMarkup {
	btnSubscribed := tele.Btn{Text: "✅ Я подписался"}
	btnPhone := tele.Btn{Contact: true, Text: "📱 Отправить номер"}

	markup := &tele.ReplyMarkup{}
	markup.Reply(
		markup.Row(btnSubscribed),
		markup.Row(btnPhone),
	)
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
