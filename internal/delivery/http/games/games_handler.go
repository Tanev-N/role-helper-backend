package games

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"

	"role-helper/internal/delivery/middleware"
	"role-helper/internal/models"
)

type createSessionResponse struct {
	Session          *models.Session  `json:"session"`
	PreviousSessions []models.Session `json:"previous_sessions"`
}

type enterSessionRequest struct {
	SessionKey  string `json:"session_key"`
	CharacterID string `json:"character_id"`
}

type createGameRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type finishSessionRequest struct {
	Summary string `json:"summary"`
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
	gameIDStr := vars["game_id"]
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный идентификатор игры")
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
	req.CharacterID = strings.TrimSpace(req.CharacterID)

	switch {
	case req.SessionKey == "":
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty session key"), "Код сессии не может быть пустым")
		return
	case req.CharacterID == "":
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty character id"), "Необходимо указать персонажа")
		return
	}

	player := &models.GamePlayer{
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
	sessionIDStr := vars["session_id"]
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный идентификатор сессии")
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

func (gr *GameRouter) GetGamePlayers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	gameIDStr := vars["game_id"]
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный идентификатор игры")
		return
	}

	players, err := gr.GameUsecase.GetGamePlayers(gameID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить игроков игры")
		return
	}

	writeSuccessResponse(w, http.StatusOK, players)
}
