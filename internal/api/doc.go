package api

import (
	"net/http"
)

// GetDocumentation возвращает документацию API
func GetDocumentation(content []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}
