package http

import (
	"time"
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
	Title           string   `json:"title" example:"Заказ на доставку"`
	Description     string   `json:"description" example:"Нужно привезти из Москвы в СПб"`
	StoreLink       *string  `json:"store_link,omitempty" example:"https://store.com/item/123"`
	Cost            *float64 `json:"cost,omitempty" example:"1000"`
	Reward          float64  `json:"reward" example:"100"`
	OriginCity      string   `json:"origin_city" example:"Москва"`
	DestinationCity string   `json:"destination_city" example:"Санкт-Петербург"`
	StartDate       string   `json:"start_date" example:"01/06/25"` // dd/mm/yy
	EndDate         string   `json:"end_date" example:"05/06/25"`   // dd/mm/yy
}

type OrderResponse struct {
	ID          int    `json:"id" example:"42"`
	OrderNumber string `json:"order_number" example:"Заказ #0123-0042"`
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
