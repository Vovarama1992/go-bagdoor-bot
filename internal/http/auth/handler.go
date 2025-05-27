package http_auth

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/auth"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
)

type AuthDeps struct {
	UserService *user.Service
}

type TelegramAuthRequest struct {
	InitData string `json:"init_data"`
}

type TelegramAuthResponse struct {
	AccessToken string `json:"access_token"`
}

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
func (deps AuthDeps) TelegramAuthHandler(w http.ResponseWriter, r *http.Request) {
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

func RegisterRoutes(mux *http.ServeMux, deps AuthDeps) {
	mux.HandleFunc("/auth/telegram", deps.TelegramAuthHandler)
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
