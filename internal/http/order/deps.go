package http_order

import (
	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/storage"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
)

type OrderDeps struct {
	OrderService *order.Service
	UserService  *user.Service
	Uploader     *storage.S3Uploader
}
