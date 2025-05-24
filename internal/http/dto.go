package http

import (
	"time"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
)

// --- Users ---

type UserResponse struct {
	ID           int       `json:"id"`
	TgID         int64     `json:"tg_id"`
	TgUsername   string    `json:"tg_username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PhoneNumber  string    `json:"phone_number"`
	RegisteredAt time.Time `json:"registered_at"`
}

// --- Orders ---

type OrderRequest struct {
	// Заголовок заказа
	Title string `json:"title" example:"Заказ на доставку"`

	// Описание задачи
	Description string `json:"description" example:"Нужно привезти документы из Москвы в Санкт-Петербург"`

	// Ссылка на магазин (только для типа store)
	StoreLink *string `json:"store_link,omitempty" example:"https://store.com/item/123"`

	// Стоимость товаров (только для типа store)
	Cost *float64 `json:"cost,omitempty" example:"1500"`

	// Депозит (только для типа personal)
	Deposit *float64 `json:"deposit,omitempty" example:"500"`

	// Вознаграждение исполнителю
	Reward float64 `json:"reward" example:"100"`

	// Город отправления
	OriginCity string `json:"origin_city" example:"Москва"`

	// Город назначения
	DestinationCity string `json:"destination_city" example:"Санкт-Петербург"`

	// Начало периода
	StartDate string `json:"start_date" example:"01/06/25"` // формат dd/mm/yy

	// Конец периода
	EndDate string `json:"end_date" example:"05/06/25"` // формат dd/mm/yy

	// Тип заказа: "personal" или "store"
	// 🔷 Для "store" обязательны поля `cost` и `store_link`
	// 🔷 Для "personal" обязателен `deposit`
	Type order.OrderType `json:"type" example:"personal"`
}

type OrderResponse struct {
	ID          int    `json:"id"`
	OrderNumber string `json:"order_number"`
}

type OrderFullResponse struct {
	ID                int       `json:"id"`
	OrderNumber       string    `json:"order_number"`
	PublisherUsername string    `json:"publisher_username"`
	PublisherTgID     int64     `json:"publisher_tg_id"`
	PublishedAt       time.Time `json:"published_at"`

	OriginCity      string    `json:"origin_city"`
	DestinationCity string    `json:"destination_city"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`

	Title       string  `json:"title"`
	Description string  `json:"description"`
	StoreLink   *string `json:"store_link,omitempty"`

	Reward    float64                `json:"reward"`
	Deposit   *float64               `json:"deposit,omitempty"`
	Cost      *float64               `json:"cost,omitempty"`
	MediaURLs []string               `json:"media_urls"`
	Type      order.OrderType        `json:"type"`
	Status    order.ModerationStatus `json:"status"`
}

type PhotoUploadResponse struct {
	Uploaded int `json:"uploaded" example:"3"`
}

// --- Telegram Auth ---

type TelegramAuthRequest struct {
	InitData string `json:"initData" example:"..."` // данные из телеги
}

type TelegramAuthResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."`
	TgID        int64  `json:"tg_id"`
}

// --- Flights ---

type FlightRequest struct {
	Description string `json:"description" example:"Лечу налегке, могу взять документы."`
	Origin      string `json:"origin" example:"Санкт-Петербург"`
	Destination string `json:"destination" example:"Москва"`
	FlightDate  string `json:"flight_date" example:"10/06/25"` // dd/mm/yy
}

type FlightResponse struct {
	ID           int    `json:"id" example:"17"`
	FlightNumber string `json:"flight_number" example:"Рейс #1234-5678"`
}

type FlightFullResponse struct {
	ID                int       `json:"id"`
	FlightNumber      string    `json:"flight_number"`
	PublisherUsername string    `json:"publisher_username"`
	PublisherTgID     int64     `json:"publisher_tg_id"`
	PublishedAt       time.Time `json:"published_at"`
	FlightDate        time.Time `json:"flight_date"`

	Description string  `json:"description"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Status      string  `json:"status"`
	MapURL      *string `json:"map_url,omitempty"`
}
