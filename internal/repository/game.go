package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"role-helper/internal/models"
)

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) models.GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) CreateGame(game *models.Game) (*models.Game, error) {
	query := `INSERT INTO games (name, master_id) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRow(query, game.Name, game.MasterID).Scan(&game.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

func (r *GameRepository) CreateSession(session *models.Session) (*models.Session, error) {
	query := `INSERT INTO sessions (game_id, session_key) 
        VALUES ($1, 
            UPPER(SUBSTRING(MD5(RANDOM()::TEXT) FROM 1 FOR 12))
        ) 
        RETURNING id, session_key`

	err := r.db.QueryRow(query, session.GameID).Scan(&session.ID, &session.SessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (r *GameRepository) GetSessionByKey(sessionKey string) (*models.Session, error) {
	var session models.Session

	query := `
		SELECT id, game_id, session_key, summary, created_at 
		FROM sessions 
		WHERE session_key = $1
	`

	err := r.db.QueryRow(query, sessionKey).Scan(
		&session.ID,
		&session.GameID,
		&session.SessionKey,
		&session.Summary,
		&session.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (r *GameRepository) GetPreviousSessions(gameID int) ([]models.Session, error) {
	query := `
		SELECT id, game_id, session_key, summary, created_at, finished_at 
		FROM sessions 
		WHERE game_id = $1 AND finished_at IS NOT NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		err := rows.Scan(
			&session.ID,
			&session.GameID,
			&session.SessionKey,
			&session.Summary,
			&session.CreatedAt,
			&session.FinishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *GameRepository) AddPlayerToGame(player *models.GamePlayer) error {
	query := `
        INSERT INTO game_players (game_id, user_id, character_id) 
        VALUES ($1, $2, $3)
    `

	err := r.db.QueryRow(query, player.GameID, player.UserID, player.CharacterID).Scan(&player.ID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	return nil
}

func (r *GameRepository) FinishSession(id int, summary string) error {
	query := `
		UPDATE sessions 
		SET summary = $1 
		WHERE id = $2
	`

	result, err := r.db.Exec(query, summary, id)
	if err != nil {
		return fmt.Errorf("failed to finish session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("failed to finish session: %w", err)
	}

	return nil
}

func (r *GameRepository) GetGameByID(gameID int) (*models.Game, error) {
	var game models.Game

	query := `SELECT id, name, master_id FROM games WHERE id = $1`

	err := r.db.QueryRow(query, gameID).Scan(
		&game.ID,
		&game.Name,
		&game.MasterID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return &game, nil
}

func (r *GameRepository) GetSessionByID(sessionID int) (*models.Session, error) {
	var session models.Session

	query := `SELECT id, game_id, session_key, summary, created_at, finished_at FROM sessions WHERE id = $1`

	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.GameID,
		&session.SessionKey,
		&session.Summary,
		&session.CreatedAt,
		&session.FinishedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}
