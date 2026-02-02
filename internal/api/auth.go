package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexl/go-fake-api/internal/models"
	"github.com/alexl/go-fake-api/internal/storage"
	"github.com/alexl/go-fake-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

// Registration обработчик регистрации
func Registration(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegistrationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request", nil)
			return
		}

		// Валидация
		if validationErrors := utils.ValidateRegistration(req); len(validationErrors) > 0 {
			utils.RespondWithValidationError(w, validationErrors)
			return
		}

		// Проверка уникальности email
		if _, err := store.GetUserByEmail(req.Email); err == nil {
			validationErrors := map[string][]string{
				"email": {"user with this email already exists"},
			}
			utils.RespondWithValidationError(w, validationErrors)
			return
		}

		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", nil)
			return
		}

		// Создание пользователя
		user := &models.User{
			Name:      req.Name,
			Email:     req.Email,
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
		}

		if err := store.CreateUser(user); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user", nil)
			return
		}

		// Ответ
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]string{
					"name":  user.Name,
					"email": user.Email,
				},
				"code":    201,
				"message": "Пользователь создан",
			},
		}

		utils.RespondWithJSON(w, http.StatusCreated, response)
	}
}

// Authorization обработчик авторизации
func Authorization(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AuthorizationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request", nil)
			return
		}

		// Валидация
		if validationErrors := utils.ValidateAuthorization(req); len(validationErrors) > 0 {
			utils.RespondWithValidationError(w, validationErrors)
			return
		}

		// Поиск пользователя
		user, err := store.GetUserByEmail(req.Email)
		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, "Login failed", nil)
			return
		}

		// Проверка пароля
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			utils.RespondWithError(w, http.StatusForbidden, "Login failed", nil)
			return
		}

		// Генерация JWT токена
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token", nil)
			return
		}

		// Сохранение токена
		if err := store.UpdateUserToken(user.ID, tokenString); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save token", nil)
			return
		}

		// Ответ
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"id":    user.ID,
					"name":  user.Name,
					"email": user.Email,
				},
				"token": tokenString,
			},
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

// Logout обработчик выхода
func Logout(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.User)
		
		// Очистка токена
		if err := store.UpdateUserToken(user.ID, ""); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to logout", nil)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
