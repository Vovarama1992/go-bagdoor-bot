package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/auth"
)

type contextKey string

const ctxKeyTgID contextKey = "tg_id"

func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		if token == "" {
			http.Error(w, "Отсутствует токен", http.StatusUnauthorized)
			return
		}

		tgID, err := auth.ParseTokenAndExtractTgID(token)
		if err != nil {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyTgID, tgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetTgIDFromContext(ctx context.Context) (int64, bool) {
	tgID, ok := ctx.Value(ctxKeyTgID).(int64)
	return tgID, ok
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}
