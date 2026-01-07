package games

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"role-helper/internal/delivery/middleware"
	"role-helper/internal/models"
)

type createSessionResponse struct {
	Session          *models.Session  `json:"session"`
	PreviousSessions []models.Session `json:"previous_sessions"`
}

type enterSessionRequest struct {
	SessionKey  string `json:"session_key"`
	CharacterID int    `json:"character_id"`
}

type createGameRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
}

type finishSessionRequest struct {
	Summary string `json:"summary"`
}

func (gr *GameRouter) LeaveSession(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	if sessionID == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty session_id"), "Некорректный идентификатор сессии")
		return
	}

	if err := gr.GameUsecase.LeaveSession(sessionID, user.ID); err != nil {
		status := http.StatusInternalServerError
		message := "Не удалось выйти из сессии"

		switch {
		case strings.Contains(err.Error(), "invalid session"):
			status = http.StatusNotFound
			message = "Сессия не найдена"
		case strings.Contains(err.Error(), "empty session id"):
			status = http.StatusBadRequest
			message = "Некорректный идентификатор сессии"
		case strings.Contains(err.Error(), "player not in session"):
			status = http.StatusBadRequest
			message = "Пользователь не состоит в сессии"
		}

		writeErrorResponse(w, status, err, message)
		return
	}

	writeSuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Пользователь вышел из сессии",
	})
}

func (gr *GameRouter) CreateGame(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	var req createGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	game := &models.Game{
		Name:        req.Name,
		MasterID:    user.ID,
		Description: req.Description,
		Photo:       req.Photo,
	}

	createdGame, err := gr.GameUsecase.CreateGame(game)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать игру")
		return
	}
	writeSuccessResponse(w, http.StatusCreated, createdGame)
}

func (gr *GameRouter) GetGames(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	games, err := gr.GameUsecase.GetAllGames(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить список игр")
		return
	}

	writeSuccessResponse(w, http.StatusOK, games)
}

func (gr *GameRouter) CreateSession(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	gameID := vars["game_id"]
	if gameID == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty game_id"), "Некорректный идентификатор игры")
		return
	}

	session, previous, err := gr.GameUsecase.CreateSession(&models.Session{GameID: gameID})
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать сессию")
		return
	}

	response := createSessionResponse{
		Session:          session,
		PreviousSessions: previous,
	}

	writeSuccessResponse(w, http.StatusCreated, response)
}

func (gr *GameRouter) EnterSession(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	var req enterSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	req.SessionKey = strings.TrimSpace(req.SessionKey)

	switch {
	case req.SessionKey == "":
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty session key"), "Код сессии не может быть пустым")
		return
	case req.CharacterID == 0:
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty character id"), "Необходимо указать персонажа")
		return
	}

	player := &models.SessionPlayer{
		UserID:      user.ID,
		CharacterID: req.CharacterID,
	}

	session, err := gr.GameUsecase.EnterSession(req.SessionKey, player)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Не удалось присоединиться к сессии"

		switch {
		case errors.Is(err, models.ErrUserNotFound):
			status = http.StatusUnauthorized
			message = "Пользователь не найден"
		case strings.Contains(err.Error(), "session not found"):
			status = http.StatusNotFound
			message = "Сессия не найдена"
		case strings.Contains(err.Error(), "invalid character") || strings.Contains(err.Error(), "invalid user id"):
			status = http.StatusBadRequest
			message = "Персонаж не принадлежит пользователю"
		case strings.Contains(err.Error(), "invalid player"):
			status = http.StatusBadRequest
			message = "Некорректный игрок"
		case strings.Contains(err.Error(), "empty session key"):
			status = http.StatusBadRequest
			message = "Код сессии не может быть пустым"
		case strings.Contains(err.Error(), "session is already finished"):
			status = http.StatusBadRequest
			message = "Сессия уже завершена"
		}

		writeErrorResponse(w, status, err, message)
		return
	}

	writeSuccessResponse(w, http.StatusOK, session)
}

func (gr *GameRouter) FinishSession(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	if sessionID == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty session_id"), "Некорректный идентификатор сессии")
		return
	}

	var req finishSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	summary := strings.TrimSpace(req.Summary)

	if err := gr.GameUsecase.FinishSession(sessionID, summary); err != nil {
		status := http.StatusInternalServerError
		message := "Не удалось завершить сессию"

		if strings.Contains(err.Error(), "invalid session") {
			status = http.StatusNotFound
			message = "Сессия не найдена"
		}

		writeErrorResponse(w, status, err, message)
		return
	}

	writeSuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Сессия завершена",
	})
}

func (gr *GameRouter) GetSessionPlayers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	if sessionID == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty session_id"), "Некорректный идентификатор сессии")
		return
	}

	players, err := gr.GameUsecase.GetActiveSessionPlayers(sessionID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Не удалось получить игроков активной сессии"

		switch {
		case strings.Contains(err.Error(), "invalid session"):
			status = http.StatusNotFound
			message = "Сессия не найдена"
		case strings.Contains(err.Error(), "empty session id"):
			status = http.StatusBadRequest
			message = "Некорректный идентификатор сессии"
		case strings.Contains(err.Error(), "session is already finished"):
			status = http.StatusBadRequest
			message = "Сессия уже завершена"
		}

		writeErrorResponse(w, status, err, message)
		return
	}

	writeSuccessResponse(w, http.StatusOK, players)
}

func (gr *GameRouter) GetPreviousSessionPlayers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	if sessionID == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty session_id"), "Некорректный идентификатор сессии")
		return
	}

	players, err := gr.GameUsecase.GetFinishedSessionPlayers(sessionID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Не удалось получить игроков завершенной сессии"

		switch {
		case strings.Contains(err.Error(), "invalid session"):
			status = http.StatusNotFound
			message = "Сессия не найдена"
		case strings.Contains(err.Error(), "empty session id"):
			status = http.StatusBadRequest
			message = "Некорректный идентификатор сессии"
		case strings.Contains(err.Error(), "session is not finished"):
			status = http.StatusBadRequest
			message = "Сессия еще не завершена"
		}

		writeErrorResponse(w, status, err, message)
		return
	}

	writeSuccessResponse(w, http.StatusOK, players)
}

func (gr *GameRouter) GetGamePlayers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	gameID := vars["game_id"]
	if gameID == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty game_id"), "Некорректный идентификатор игры")
		return
	}

	players, err := gr.GameUsecase.GetGamePlayers(gameID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить игроков игры")
		return
	}

	writeSuccessResponse(w, http.StatusOK, players)
}

func (gr *GameRouter) GetPreviousSessions(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	gameID := vars["game_id"]
	if gameID == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty game_id"), "Некорректный идентификатор игры")
		return
	}

	sessions, err := gr.GameUsecase.GetPreviousSessions(gameID, user.ID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Не удалось получить предыдущие сессии"

		switch {
		case strings.Contains(err.Error(), "game not found"):
			status = http.StatusNotFound
			message = "Игра не найдена"
		case strings.Contains(err.Error(), "empty game id"):
			status = http.StatusBadRequest
			message = "Некорректный идентификатор игры"
		case strings.Contains(err.Error(), "access denied"):
			status = http.StatusForbidden
			message = "Нет доступа к предыдущим сессиям"
		}

		writeErrorResponse(w, status, err, message)
		return
	}

	writeSuccessResponse(w, http.StatusOK, sessions)
}
