package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexl/go-fake-api/internal/models"
	"github.com/alexl/go-fake-api/internal/storage"
	"github.com/alexl/go-fake-api/internal/utils"
	"github.com/gorilla/mux"
)

// CreateBoard создает новую доску
func CreateBoard(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.User)

		var req models.BoardCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.SendError(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}

		boardID := fmt.Sprintf("board-%d", time.Now().UnixNano())
		hash := fmt.Sprintf("%x", time.Now().UnixNano())

		board := &models.Board{
			ID:        boardID,
			Hash:      hash,
			Name:      req.Name,
			OwnerID:   user.ID,
			IsPublic:  req.IsPublic,
			Objects:   make(map[string]models.BoardObject),
			CreatedAt: time.Now(),
		}

		if err := s.CreateBoard(board); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "could not create board", nil)
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "board created", board)
	}
}

// GetUserBoards возвращает доски, к которым у пользователя есть доступ
func GetUserBoards(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.User)

		boards, err := s.GetUserBoards(user.ID)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "could not fetch boards", nil)
			return
		}

		utils.SendSuccess(w, http.StatusOK, "success", boards)
	}
}

// GetPublicBoards возвращает список публичных досок
func GetPublicBoards(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		boards, err := s.GetPublicBoards()
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "could not fetch public boards", nil)
			return
		}

		utils.SendSuccess(w, http.StatusOK, "success", boards)
	}
}

// ShareBoard предоставляет доступ к доске
func ShareBoard(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.User)
		vars := mux.Vars(r)
		boardID := vars["board_id"]

		board, err := s.GetBoardByID(boardID)
		if err != nil {
			utils.SendError(w, http.StatusNotFound, "board not found", nil)
			return
		}

		if board.OwnerID != user.ID {
			utils.SendError(w, http.StatusForbidden, "only owner can share board", nil)
			return
		}

		var req models.BoardShareRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.SendError(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}

		recipient, err := s.GetUserByEmail(req.Email)
		if err != nil {
			utils.SendError(w, http.StatusNotFound, "user with this email not found", nil)
			return
		}

		if err := s.AddBoardAccess(boardID, recipient.ID); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "could not share board", nil)
			return
		}

		utils.SendSuccess(w, http.StatusOK, "board shared successfully", nil)
	}
}

// GetBoardByHash возвращает публичную информацию о доске
func GetBoardByHash(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]

		board, err := s.GetBoardByHash(hash)
		if err != nil {
			utils.SendError(w, http.StatusNotFound, "board not found", nil)
			return
		}

		utils.SendSuccess(w, http.StatusOK, "success", board)
	}
}

// LikeBoard ставит лайк доске
func LikeBoard(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.User)
		vars := mux.Vars(r)
		boardID := vars["board_id"]

		if err := s.LikeBoard(boardID, user.ID); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "could not process like", nil)
			return
		}

		utils.SendSuccess(w, http.StatusOK, "success", nil)
	}
}
