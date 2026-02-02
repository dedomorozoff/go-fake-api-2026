package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alexl/go-fake-api/internal/api"
	"github.com/alexl/go-fake-api/internal/middleware"
	"github.com/alexl/go-fake-api/internal/storage"
	"github.com/gorilla/mux"
	_ "embed"
)

//go:embed API_DOCUMENTATION.md
var documentation []byte

func main() {
	// Инициализация хранилища
	store := storage.NewMemoryStorage()

	// Инициализация Hub для WebSocket
	hub := api.NewHub(store)
	go hub.Run()

	// Создание роутера
	r := mux.NewRouter()

	// Middleware для CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Публичные эндпоинты
	r.HandleFunc("/", api.GetDocumentation(documentation)).Methods("GET")
	r.HandleFunc("/registration", api.Registration(store)).Methods("POST")
	r.HandleFunc("/authorization", api.Authorization(store)).Methods("POST")
	r.HandleFunc("/public-boards", api.GetPublicBoards(store)).Methods("GET")
	r.HandleFunc("/board/{hash}", api.GetBoardByHash(store)).Methods("GET")

	// Защищенные эндпоинты
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(store))

	protected.HandleFunc("/logout", api.Logout(store)).Methods("GET")
	protected.HandleFunc("/boards", api.CreateBoard(store)).Methods("POST")
	protected.HandleFunc("/boards", api.GetUserBoards(store)).Methods("GET")
	protected.HandleFunc("/boards/{board_id}/share", api.ShareBoard(store)).Methods("POST")
	protected.HandleFunc("/boards/{board_id}/like", api.LikeBoard(store)).Methods("POST")

	// WebSocket
	r.HandleFunc("/ws/board/{board_id}", api.ServeWs(hub, store))

	// Получение порта из переменной окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
