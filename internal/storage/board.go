package storage

import (
	"errors"
	"sort"

	"github.com/alexl/go-fake-api/internal/models"
)

// CreateBoard создает новую доску
func (s *MemoryStorage) CreateBoard(board *models.Board) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.boards[board.ID] = board
	s.boardAccess[board.ID] = append(s.boardAccess[board.ID], board.OwnerID)
	s.boardLikes[board.ID] = make(map[int]bool)
	return nil
}

// GetBoardByID возвращает доску по ID
func (s *MemoryStorage) GetBoardByID(id string) (*models.Board, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	board, ok := s.boards[id]
	if !ok {
		return nil, errors.New("board not found")
	}
	return board, nil
}

// GetBoardByHash возвращает доску по хешу
func (s *MemoryStorage) GetBoardByHash(hash string) (*models.Board, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, board := range s.boards {
		if board.Hash == hash {
			return board, nil
		}
	}
	return nil, errors.New("board not found")
}

// GetUserBoards возвращает список досок пользователя
func (s *MemoryStorage) GetUserBoards(userID int) ([]models.Board, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userBoards []models.Board
	for boardID, userIDs := range s.boardAccess {
		for _, uid := range userIDs {
			if uid == userID {
				if board, ok := s.boards[boardID]; ok {
					userBoards = append(userBoards, *board)
				}
				break
			}
		}
	}
	return userBoards, nil
}

// GetPublicBoards возвращает список публичных досок
func (s *MemoryStorage) GetPublicBoards() ([]models.Board, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var publicBoards []models.Board
	for _, board := range s.boards {
		if board.IsPublic {
			publicBoards = append(publicBoards, *board)
		}
	}

	// Сортировка по количеству лайков (убывание)
	sort.Slice(publicBoards, func(i, j int) bool {
		return publicBoards[i].Likes > publicBoards[j].Likes
	})

	return publicBoards, nil
}

// UpdateBoardObject обновление или добавление объекта на доске
func (s *MemoryStorage) UpdateBoardObject(boardID string, obj models.BoardObject) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	board, ok := s.boards[boardID]
	if !ok {
		return errors.New("board not found")
	}

	if board.Objects == nil {
		board.Objects = make(map[string]models.BoardObject)
	}
	board.Objects[obj.ID] = obj
	return nil
}

// DeleteBoardObject удаляет объект с доски
func (s *MemoryStorage) DeleteBoardObject(boardID string, objectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	board, ok := s.boards[boardID]
	if !ok {
		return errors.New("board not found")
	}

	delete(board.Objects, objectID)
	return nil
}

// AddBoardAccess предоставляет доступ к доске
func (s *MemoryStorage) AddBoardAccess(boardID string, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, есть ли уже доступ
	for _, uid := range s.boardAccess[boardID] {
		if uid == userID {
			return nil
		}
	}

	s.boardAccess[boardID] = append(s.boardAccess[boardID], userID)
	return nil
}

// HasBoardAccess проверяет наличие доступа
func (s *MemoryStorage) HasBoardAccess(boardID string, userID int) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userIDs, ok := s.boardAccess[boardID]
	if !ok {
		return false, nil
	}

	for _, uid := range userIDs {
		if uid == userID {
			return true, nil
		}
	}
	return false, nil
}

// LikeBoard ставит/снимает лайк
func (s *MemoryStorage) LikeBoard(boardID string, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	board, ok := s.boards[boardID]
	if !ok {
		return errors.New("board not found")
	}

	if s.boardLikes[boardID] == nil {
		s.boardLikes[boardID] = make(map[int]bool)
	}

	if s.boardLikes[boardID][userID] {
		// Убираем лайк
		delete(s.boardLikes[boardID], userID)
		board.Likes--
	} else {
		// Ставим лайк
		s.boardLikes[boardID][userID] = true
		board.Likes++
	}

	return nil
}
