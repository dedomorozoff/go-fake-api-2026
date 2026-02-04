package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexl/go-fake-api/internal/api"
	"github.com/alexl/go-fake-api/internal/middleware"
	"github.com/alexl/go-fake-api/internal/storage"
	"github.com/gorilla/mux"
	_ "embed"
)

//go:embed API_DOCUMENTATION.md
var documentation []byte

func main() {
	// Парсинг аргументов командной строки
	var baseURL string
	var port string
	flag.StringVar(&baseURL, "base-url", "", "Base URL path for the API (e.g., /api/v1)")
	flag.StringVar(&port, "port", "", "Port to listen on (default: 8080 or PORT env var)")
	flag.Parse()

	// Нормализация base URL
	if baseURL != "" {
		baseURL = strings.TrimSuffix(baseURL, "/")
		if !strings.HasPrefix(baseURL, "/") {
			baseURL = "/" + baseURL
		}
	}

	// Инициализация хранилища
	store := storage.NewMemoryStorage()

	// Инициализация Hub для WebSocket
	hub := api.NewHub(store)
	go hub.Run()

	// Создание роутера
	r := mux.NewRouter()

	// Middleware для CORS - полная отмена проверок для тестового API
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Разрешаем все origins
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			
			// Разрешаем все методы и заголовки
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")
			
			// Обрабатываем preflight запросы
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	// Создание подроутера с базовым URL если указан
	var apiRouter *mux.Router
	if baseURL != "" {
		apiRouter = r.PathPrefix(baseURL).Subrouter()
	} else {
		apiRouter = r
	}

	// Публичные эндпоинты
	apiRouter.HandleFunc("/", api.GetDocumentation(documentation)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/registration", api.Registration(store)).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/authorization", api.Authorization(store)).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/public-boards", api.GetPublicBoards(store)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/board/{hash}", api.GetBoardByHash(store)).Methods("GET", "OPTIONS")

	// Защищенные эндпоинты
	protected := apiRouter.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(store))

	// Добавляем OPTIONS методы для всех защищенных эндпоинтов
	protected.HandleFunc("/logout", api.Logout(store)).Methods("GET", "OPTIONS")
	protected.HandleFunc("/boards", api.CreateBoard(store)).Methods("POST", "OPTIONS")
	protected.HandleFunc("/boards", api.GetUserBoards(store)).Methods("GET", "OPTIONS")
	protected.HandleFunc("/boards/{board_id}/share", api.ShareBoard(store)).Methods("POST", "OPTIONS")
	protected.HandleFunc("/boards/{board_id}/like", api.LikeBoard(store)).Methods("POST", "OPTIONS")

	// WebSocket
	apiRouter.HandleFunc("/ws/board/{board_id}", api.ServeWs(hub, store))

	// Получение порта из аргумента командной строки или переменной окружения
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	if baseURL != "" {
		log.Printf("Server starting on port %s with base URL %s...", port, baseURL)
	} else {
		log.Printf("Server starting on port %s...", port)
	}
	
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
