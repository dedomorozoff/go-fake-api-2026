package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/alexl/go-fake-api/internal/models"
)

// ValidateRegistration валидирует данные регистрации
func ValidateRegistration(req models.RegistrationRequest) map[string][]string {
	errors := make(map[string][]string)

	// Проверка name (только латиница)
	if req.Name == "" {
		errors["name"] = append(errors["name"], "field name can not be blank")
	} else if !isLatin(req.Name) {
		errors["name"] = append(errors["name"], "name must contain only latin characters")
	}

	// Проверка email
	if req.Email == "" {
		errors["email"] = append(errors["email"], "field email can not be blank")
	} else if !isValidEmail(req.Email) {
		errors["email"] = append(errors["email"], "invalid email format")
	}

	// Проверка password (от 8 символов, цифры и спецсимволы)
	if req.Password == "" {
		errors["password"] = append(errors["password"], "field password can not be blank")
	} else if len(req.Password) < 8 {
		errors["password"] = append(errors["password"], "password must be at least 8 characters long")
	} else if !hasDigitsAndSpecials(req.Password) {
		errors["password"] = append(errors["password"], "password must contain digits and special characters")
	}

	return errors
}

// ValidateAuthorization валидирует данные авторизации
func ValidateAuthorization(req models.AuthorizationRequest) map[string][]string {
	errors := make(map[string][]string)

	if req.Email == "" {
		errors["email"] = append(errors["email"], "field email can not be blank")
	}

	if req.Password == "" {
		errors["password"] = append(errors["password"], "field password can not be blank")
	}

	return errors
}

// isLatin проверяет, содержит ли строка только латинские буквы
func isLatin(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && r != ' ' {
			return false
		}
	}
	return true
}

// hasDigitsAndSpecials проверяет наличие цифр и специальных символов
func hasDigitsAndSpecials(s string) bool {
	hasDigit := false
	hasSpecial := false

	for _, r := range s {
		if unicode.IsDigit(r) {
			hasDigit = true
		} else if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			hasSpecial = true
		}
	}

	return hasDigit && hasSpecial
}

// isValidEmail проверяет валидность email
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// NormalizeString нормализует строку (первая буква заглавная, остальные строчные)
func NormalizeString(s string) string {
	if len(s) == 0 {
		return s
	}
	
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
