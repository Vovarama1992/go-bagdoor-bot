package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/auth"
)

// @Summary Создать рейс
// @Tags flights
// @Accept json
// @Produce json
// @Param flight body FlightRequest true "Данные рейса"
// @Success 201 {object} FlightResponse
// @Failure 400 {string} string "Невалидный JSON или формат даты"
// @Failure 401 {string} string "Неверный или отсутствует токен"
// @Failure 500 {string} string "Ошибка сервера при создании рейса"
// @Router /flights [post]
func (d FlightDeps) createFlightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Нет токена", http.StatusUnauthorized)
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	tgID, err := auth.ParseTokenAndExtractTgID(token)
	if err != nil {
		http.Error(w, "Невалидный токен", http.StatusUnauthorized)
		return
	}

	var req FlightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	flightDate, err := time.Parse("02/01/06", req.FlightDate)
	if err != nil {
		http.Error(w, "Невалидная дата, используйте формат dd/mm/yy", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u, err := d.UserService.GetByTgID(ctx, tgID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusInternalServerError)
		return
	}

	f, err := d.FlightService.CreateFlight(
		ctx,
		u.TgUsername,
		tgID,
		req.Description,
		req.Origin,
		req.Destination,
		flightDate,
	)
	if err != nil {
		http.Error(w, "Не удалось создать рейс", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(FlightResponse{
		ID:           f.ID,
		FlightNumber: f.FlightNumber,
	})
}

func RegisterFlightRoutes(mux *http.ServeMux, deps FlightDeps) {
	mux.HandleFunc("/flights", WithAuth(deps.createFlightHandler))
}
