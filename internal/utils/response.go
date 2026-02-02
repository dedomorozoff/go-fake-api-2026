package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse структура ответа с ошибкой
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// ValidationError структура ошибки валидации
type ValidationError struct {
	Error ValidationErrorDetail `json:"error"`
}

// ValidationErrorDetail детали ошибки валидации
type ValidationErrorDetail struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Errors  map[string][]string       `json:"errors"`
}

// RespondWithJSON отправляет JSON ответ
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// RespondWithError отправляет ответ с ошибкой
func RespondWithError(w http.ResponseWriter, statusCode int, message string, code *int) {
	response := ErrorResponse{
		Message: message,
	}
	
	if code != nil {
		response.Code = *code
	}
	
	RespondWithJSON(w, statusCode, response)
}

// RespondWithValidationError отправляет ответ с ошибкой валидации
func RespondWithValidationError(w http.ResponseWriter, errors map[string][]string) {
	response := ValidationError{
		Error: ValidationErrorDetail{
			Code:    422,
			Message: "Validation error",
			Errors:  errors,
		},
	}
	
	RespondWithJSON(w, http.StatusUnprocessableEntity, response)
}

// SendSuccess отправляет успешный JSON ответ
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := map[string]interface{}{
		"data": data,
	}
	if message != "" {
		response["message"] = message
	}
	RespondWithJSON(w, statusCode, response)
}

// SendError отправляет ответ с ошибкой
func SendError(w http.ResponseWriter, statusCode int, message string, errors map[string][]string) {
	if errors != nil {
		RespondWithValidationError(w, errors)
		return
	}
	RespondWithError(w, statusCode, message, nil)
}
