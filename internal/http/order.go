package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/auth"
	telegram "github.com/Vovarama1992/go-bagdoor-bot/internal/notifier"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
)

// @Summary Авторизация через Telegram
// @Tags auth
// @Accept json
// @Produce json
// @Param data body TelegramAuthRequest true "Init data от Telegram"
// @Success 200 {object} TelegramAuthResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Invalid signature"
// @Failure 500 {string} string "User error или Token error"
// @Router /auth/telegram [post]
func (deps OrderDeps) telegramAuthHandler(w http.ResponseWriter, r *http.Request) {
	var body TelegramAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	data, ok := auth.ValidateTelegramInitData(body.InitData, os.Getenv("TELEGRAM_BOT_TOKEN"))
	if !ok {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	tgID, _ := parseInt64(data["user.id"])
	username := data["user.username"]
	first := data["user.first_name"]
	last := data["user.last_name"]

	ctx := r.Context()
	if _, err := deps.UserService.FindOrCreateFromTelegram(ctx, tgID, username, first, last); err != nil {
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateTokenWithTgID(tgID)
	if err != nil {
		http.Error(w, "Token error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TelegramAuthResponse{
		AccessToken: token,
	})
}

// @Summary Создать заказ
// @Tags orders
// @Accept json
// @Produce json
// @Param order body OrderRequest true "Данные заказа"
// @Success 201 {object} OrderResponse
// @Failure 400 {string} string "Невалидный JSON или формат дат"
// @Failure 401 {string} string "Невалидный токен"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {string} string "Ошибка сервера при создании заказа"
// @Router /orders [post]
func (deps OrderDeps) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	tgID, ok := GetTgIDFromContext(r.Context())
	if !ok {
		http.Error(w, "tgID не найден в контексте", http.StatusUnauthorized)
		return
	}

	var body OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	start, err1 := time.Parse("02/01/06", body.StartDate)
	end, err2 := time.Parse("02/01/06", body.EndDate)
	if err1 != nil || err2 != nil {
		http.Error(w, "Неверный формат дат", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := deps.UserService.GetByTgID(ctx, tgID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	o := &order.Order{
		OrderNumber:       fmt.Sprintf("Заказ #%04d-%04d", time.Now().Unix()%10000, user.ID%10000),
		PublisherUsername: user.TgUsername,
		PublisherTgID:     tgID,
		PublishedAt:       time.Now(),
		UserID:            int(user.ID),
		Title:             body.Title,
		Description:       body.Description,
		StoreLink:         body.StoreLink,
		Cost:              body.Cost,
		Reward:            body.Reward,
		OriginCity:        body.OriginCity,
		DestinationCity:   body.DestinationCity,
		StartDate:         start,
		EndDate:           end,
		Type:              order.OrderTypePersonal,
		Status:            order.StatusPending,
	}

	if err := deps.OrderService.CreateOrder(ctx, o); err != nil {
		http.Error(w, "Ошибка при создании заказа", http.StatusInternalServerError)
		return
	}

	notifier := telegram.NewNotifier(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err := notifier.NotifyNewOrder(tgID, o); err != nil {
		log.Printf("Ошибка отправки Telegram-уведомления: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(OrderResponse{
		ID:          o.ID,
		OrderNumber: o.OrderNumber,
	})
}

// RegisterRoutes регистрирует роуты
func RegisterRoutes(mux *http.ServeMux, deps OrderDeps) {
	mux.HandleFunc("/orders", WithAuth(deps.createOrderHandler))
	mux.HandleFunc("/auth/telegram", deps.telegramAuthHandler)
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
