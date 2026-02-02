package models

import (
	"time"
)

// User представляет пользователя системы
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Не отдаем пароль в JSON
	Token     string    `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// RegistrationRequest структура запроса регистрации
type RegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthorizationRequest структура запроса авторизации
type AuthorizationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
