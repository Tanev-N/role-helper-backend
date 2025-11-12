package models

import (
	"time"
)

type Game struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MasterID    int    `json:"master_id"`
	Description string `json:"description"`
}

type Session struct {
	ID         string     `json:"id"`
	GameID     string     `json:"game_id"`
	SessionKey string     `json:"session_key"`
	Summary    *string    `json:"summary,omitempty"`
	CreatedAt  *time.Time `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type GamePlayer struct {
	ID          int    `json:"id"`
	GameID      string `json:"game_id"`
	UserID      int    `json:"user_id"`
	CharacterID int    `json:"character_id"`
}

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
	AddPlayerToGame(player *GamePlayer) error
	FinishSession(id string, summary string) error
	GetGameByID(gameID string) (*Game, error)
	GetSessionByID(sessionID string) (*Session, error)
	GetAllGames(userID int) ([]Game, error)
	GetGamePlayers(gameID string) ([]GamePlayer, error)
}

type GameService interface {
	CreateGame(game *Game) (*Game, error)
	CreateSession(session *Session) (*Session, []Session, error)
	EnterSession(key string, player *GamePlayer) (*Session, error)
	FinishSession(id string, summary string) error
	GetAllGames(userID int) ([]Game, error)
	GetGamePlayers(gameID string) ([]GamePlayer, error)
}
