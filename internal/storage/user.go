package storage

import (
	"errors"

	"github.com/alexl/go-fake-api/internal/models"
)

// CreateUser создает нового пользователя
func (s *MemoryStorage) CreateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersByEmail[user.Email]; exists {
		return errors.New("user with this email already exists")
	}

	user.ID = s.userIDCounter
	s.userIDCounter++

	s.users[user.ID] = user
	s.usersByEmail[user.Email] = user

	return nil
}

// GetUserByEmail получает пользователя по email
func (s *MemoryStorage) GetUserByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.usersByEmail[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByToken получает пользователя по токену
func (s *MemoryStorage) GetUserByToken(token string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.usersByToken[token]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUserToken обновляет токен пользователя
func (s *MemoryStorage) UpdateUserToken(userID int, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	// Удаляем старый токен
	if user.Token != "" {
		delete(s.usersByToken, user.Token)
	}

	// Устанавливаем новый токен
	user.Token = token
	if token != "" {
		s.usersByToken[token] = user
	}

	return nil
}
