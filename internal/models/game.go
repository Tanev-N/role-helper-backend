package models

import (
	"time"
)

type Game struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MasterID    int    `json:"master_id"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
}

type Session struct {
	ID         string     `json:"id"`
	GameID     string     `json:"game_id"`
	SessionKey string     `json:"session_key"`
	Summary    *string    `json:"summary,omitempty"`
	CreatedAt  *time.Time `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type SessionPlayer struct {
	ID          int    `json:"id"`
	SessionID   string `json:"session_id"`
	UserID      int    `json:"user_id"`
	CharacterID int    `json:"character_id"`
}

type GamePlayer = SessionPlayer

type EnterSession struct {
	ID          int    `json:"id"`
	GameID      int    `json:"game_id"`
	UserID      int    `json:"user_id"`
	CharacterID int    `json:"character_id"`
	SessionKey  string `json:"session_key"`
}

type GameRepository interface {
	CreateGame(game *Game) (*Game, error)
	CreateSession(session *Session) (*Session, error)
	GetSessionByKey(sessionKey string) (*Session, error)
	GetPreviousSessions(gameID string) ([]Session, error)
	AddPlayerToSession(player *SessionPlayer) error
	FinishSession(id string, summary string) error
	GetGameByID(gameID string) (*Game, error)
	GetSessionByID(sessionID string) (*Session, error)
	GetAllGames(userID int) ([]Game, error)
	GetSessionPlayers(sessionID string) ([]SessionPlayer, error)
	GetGamePlayers(gameID string) ([]SessionPlayer, error)
}

type GameService interface {
	CreateGame(game *Game) (*Game, error)
	CreateSession(session *Session) (*Session, []Session, error)
	EnterSession(key string, player *SessionPlayer) (*Session, error)
	FinishSession(id string, summary string) error
	GetAllGames(userID int) ([]Game, error)
	// Новый метод
	GetSessionPlayers(sessionID string) ([]SessionPlayer, error)
	// Старый метод для обратной совместимости
	GetGamePlayers(gameID string) ([]SessionPlayer, error)
	GetPreviousSessions(gameID string, userID int) ([]Session, error)
}
