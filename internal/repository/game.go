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
	query := `INSERT INTO games (name, master_id, description) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(query, game.Name, game.MasterID, game.Description).Scan(&game.ID)
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
        RETURNING id, session_key, created_at`

	err := r.db.QueryRow(query, session.GameID).Scan(&session.ID, &session.SessionKey, &session.CreatedAt)
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
        RETURNING id
    `

	err := r.db.QueryRow(query, player.GameID, player.UserID, player.CharacterID).Scan(&player.ID)
	if err != nil {
		return fmt.Errorf("failed to add player to game: %w", err)
	}
	return nil
}

func (r *GameRepository) FinishSession(id int, summary string) error {
	query := `
		UPDATE sessions 
		SET summary = $1, finished_at = NOW() 
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

	query := `SELECT id, name, master_id, description FROM games WHERE id = $1`

	err := r.db.QueryRow(query, gameID).Scan(
		&game.ID,
		&game.Name,
		&game.MasterID,
		&game.Description,
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

func (r *GameRepository) GetAllGames(userID int) ([]models.Game, error) {
	query := `
		SELECT g.id, g.name, g.master_id, g.description
		FROM games g
		LEFT JOIN game_players gp ON g.id = gp.game_id AND gp.user_id = $1
		WHERE g.master_id = $1 OR gp.user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query games: %w", err)
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var game models.Game
		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.MasterID,
			&game.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}
		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating games: %w", err)
	}

	return games, nil
}

func (r *GameRepository) GetGamePlayers(gameID int) ([]models.GamePlayer, error) {
	query := `
		SELECT id, game_id, user_id, character_id
		FROM game_players 
		WHERE game_id = $1
	`

	rows, err := r.db.Query(query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to query game players: %w", err)
	}
	defer rows.Close()

	var players []models.GamePlayer
	for rows.Next() {
		var player models.GamePlayer
		err := rows.Scan(
			&player.ID,
			&player.GameID,
			&player.UserID,
			&player.CharacterID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game player: %w", err)
		}
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating game players: %w", err)
	}

	return players, nil
}
