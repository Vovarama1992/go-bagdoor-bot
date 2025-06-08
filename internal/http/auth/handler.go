package http_auth

import (
	"encoding/json"
	"fmt"
	"io"
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
// @Param data body object true "Init data от Telegram"
// @Success 200 {object} TelegramAuthResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Invalid signature"
// @Failure 500 {string} string "User error или Token error"
// @Router /auth/telegram [post]
func (deps AuthDeps) TelegramAuthHandler(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("[auth] ❌ Failed to read request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		fmt.Println("[auth] ❌ Failed to parse JSON:", err)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	var initData string
	if err := json.Unmarshal(rawMap["init_data"], &initData); err != nil {
		fmt.Println("[auth] ❌ Failed to extract init_data:", err)
		http.Error(w, "Invalid init_data", http.StatusBadRequest)
		return
	}

	fmt.Println("[auth] ✅ Received init_data:", initData)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		fmt.Println("[auth] ❌ Missing TELEGRAM_BOT_TOKEN")
		http.Error(w, "Server config error", http.StatusInternalServerError)
		return
	}

	data, ok := auth.ValidateTelegramInitData(initData, botToken)
	if !ok {
		fmt.Println("[auth] ❌ Invalid signature for init_data")
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	fmt.Println("[auth] ✅ Signature valid, extracted data:", data)

	tgID, err := parseInt64(data["user.id"])
	if err != nil {
		fmt.Println("[auth] ❌ Failed to parse user.id:", err)
		http.Error(w, "Bad Telegram ID", http.StatusBadRequest)
		return
	}

	username := data["user.username"]
	first := data["user.first_name"]
	last := data["user.last_name"]

	ctx := r.Context()
	if _, err := deps.UserService.FindOrCreateFromTelegram(ctx, tgID, username, first, last); err != nil {
		fmt.Println("[auth] ❌ Failed to create/find user:", err)
		http.Error(w, "User error", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateTokenWithTgID(tgID)
	if err != nil {
		fmt.Println("[auth] ❌ Failed to generate token:", err)
		http.Error(w, "Token error", http.StatusInternalServerError)
		return
	}

	fmt.Println("[auth] ✅ Auth successful for user:", tgID)
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
