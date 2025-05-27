package http_order

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"time"

	telegram "github.com/Vovarama1992/go-bagdoor-bot/internal/notifier"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/server"
)

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

	tgID, ok := server.GetTgIDFromContext(r.Context())
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

	switch body.Type {
	case order.OrderTypeStore:
		if body.Cost == nil || body.StoreLink == nil {
			http.Error(w, "Для типа store обязательны cost и store_link", http.StatusBadRequest)
			return
		}
	case order.OrderTypePersonal:
		if body.Deposit == nil {
			http.Error(w, "Для типа personal обязателен deposit", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Неверный тип заказа", http.StatusBadRequest)
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
		Deposit:           body.Deposit,
		Reward:            body.Reward,
		OriginCity:        body.OriginCity,
		DestinationCity:   body.DestinationCity,
		StartDate:         start,
		EndDate:           end,
		Type:              body.Type,
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

// @Summary Получить все заказы
// @Tags orders
// @Produce json
// @Success 200 {array} OrderFullResponse
// @Failure 500 {string} string "Ошибка сервера при получении заказов"
// @Router /orders [get]
func (deps OrderDeps) getAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	orders, err := deps.OrderService.GetAllOrders(ctx)
	if err != nil {
		http.Error(w, "Ошибка при получении заказов", http.StatusInternalServerError)
		return
	}

	var response []OrderFullResponse
	for _, o := range orders {
		response = append(response, OrderFullResponse{
			ID:                o.ID,
			OrderNumber:       o.OrderNumber,
			PublisherUsername: o.PublisherUsername,
			PublisherTgID:     o.PublisherTgID,
			PublishedAt:       o.PublishedAt,
			OriginCity:        o.OriginCity,
			DestinationCity:   o.DestinationCity,
			StartDate:         o.StartDate,
			EndDate:           o.EndDate,
			Title:             o.Title,
			Description:       o.Description,
			StoreLink:         o.StoreLink,
			Reward:            o.Reward,
			Deposit:           o.Deposit,
			Cost:              o.Cost,
			MediaURLs:         o.MediaURLs,
			Type:              o.Type,
			Status:            o.Status,
		})
	}

	json.NewEncoder(w).Encode(response)
}

func RegisterRoutes(mux *http.ServeMux, deps OrderDeps) {
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			server.WithAuth(deps.createOrderHandler)(w, r)
		case http.MethodGet:
			deps.getAllOrdersHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

}
