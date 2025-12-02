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

	createdSession := &models.Session{
		GameID: session.GameID,
	}

	err := r.db.QueryRow(query, session.GameID).Scan(&createdSession.ID, &createdSession.SessionKey, &createdSession.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return createdSession, nil
}

func (r *GameRepository) GetSessionByKey(sessionKey string) (*models.Session, error) {
	var session models.Session

	query := `
		SELECT id, game_id, session_key, summary, created_at, finished_at
		FROM sessions 
		WHERE session_key = $1
	`

	err := r.db.QueryRow(query, sessionKey).Scan(
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

func (r *GameRepository) GetPreviousSessions(gameID string) ([]models.Session, error) {
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

func (r *GameRepository) AddPlayerToSession(player *models.SessionPlayer) error {
	checkQuery := `SELECT id FROM session_players WHERE session_id = $1 AND user_id = $2`
	var existingID int
	err := r.db.QueryRow(checkQuery, player.SessionID, player.UserID).Scan(&existingID)
	if err == nil {
		player.ID = existingID
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check player in session: %w", err)
	}

	query := `
        INSERT INTO session_players (session_id, user_id, character_id) 
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err = r.db.QueryRow(query, player.SessionID, player.UserID, player.CharacterID).Scan(&player.ID)
	if err != nil {
		return fmt.Errorf("failed to add player to session: %w", err)
	}
	return nil
}

func (r *GameRepository) FinishSession(id string, summary string) error {
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

func (r *GameRepository) GetGameByID(gameID string) (*models.Game, error) {
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

func (r *GameRepository) GetSessionByID(sessionID string) (*models.Session, error) {
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
		WHERE g.master_id = $1
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

func (r *GameRepository) GetSessionPlayers(sessionID string) ([]models.SessionPlayer, error) {
	query := `
		SELECT id, session_id, user_id, character_id
		FROM session_players 
		WHERE session_id = $1
	`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query session players: %w", err)
	}
	defer rows.Close()

	var players []models.SessionPlayer
	for rows.Next() {
		var player models.SessionPlayer
		err := rows.Scan(
			&player.ID,
			&player.SessionID,
			&player.UserID,
			&player.CharacterID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session player: %w", err)
		}
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session players: %w", err)
	}

	return players, nil
}

func (r *GameRepository) GetGamePlayers(gameID string) ([]models.SessionPlayer, error) {
	query := `
		SELECT sp.id, sp.session_id, sp.user_id, sp.character_id
		FROM session_players sp
		JOIN sessions s ON sp.session_id = s.id
		WHERE s.game_id = $1
	`

	rows, err := r.db.Query(query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to query game players: %w", err)
	}
	defer rows.Close()

	var players []models.SessionPlayer
	for rows.Next() {
		var player models.SessionPlayer
		err := rows.Scan(
			&player.ID,
			&player.SessionID,
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
