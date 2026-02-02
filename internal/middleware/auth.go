package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexl/go-fake-api/internal/storage"
	"github.com/alexl/go-fake-api/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware проверяет Bearer токен
func AuthMiddleware(store storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.RespondWithError(w, http.StatusForbidden, "Login failed", nil)
				return
			}

			// Проверяем формат Bearer токена
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.RespondWithError(w, http.StatusForbidden, "Login failed", nil)
				return
			}

			token := parts[1]

			// Получаем пользователя по токену
			user, err := store.GetUserByToken(token)
			if err != nil {
				utils.RespondWithError(w, http.StatusForbidden, "Login failed", nil)
				return
			}

			// Добавляем пользователя в контекст
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
