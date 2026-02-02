package storage

import (
	"sync"

	"github.com/alexl/go-fake-api/internal/models"
)

// Storage интерфейс хранилища
type Storage interface {
	// Users
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByToken(token string) (*models.User, error)
	UpdateUserToken(userID int, token string) error

	// Boards
	CreateBoard(board *models.Board) error
	GetBoardByID(id string) (*models.Board, error)
	GetBoardByHash(hash string) (*models.Board, error)
	GetUserBoards(userID int) ([]models.Board, error)
	GetPublicBoards() ([]models.Board, error)
	UpdateBoardObject(boardID string, obj models.BoardObject) error
	DeleteBoardObject(boardID string, objectID string) error
	AddBoardAccess(boardID string, userID int) error
	HasBoardAccess(boardID string, userID int) (bool, error)
	LikeBoard(boardID string, userID int) error
}

// MemoryStorage хранилище в памяти
type MemoryStorage struct {
	users         map[int]*models.User
	usersByEmail  map[string]*models.User
	usersByToken  map[string]*models.User
	boards        map[string]*models.Board
	boardAccess   map[string][]int        // boardID -> []userID
	boardLikes    map[string]map[int]bool // boardID -> userID -> true
	userIDCounter int
	mu            sync.RWMutex
}

// NewMemoryStorage создает новое хранилище в памяти
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:         make(map[int]*models.User),
		usersByEmail:  make(map[string]*models.User),
		usersByToken:  make(map[string]*models.User),
		boards:        make(map[string]*models.Board),
		boardAccess:   make(map[string][]int),
		boardLikes:    make(map[string]map[int]bool),
		userIDCounter: 1,
	}
}
