package repository

import (
	"database/sql"
	"role-helper/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, avatar_url)
		VALUES ($1, $2, $3)
		RETURNING id, username, avatar_url
	`

	err := ur.db.QueryRow(query, user.Username, user.PasswordHash, user.AvatarURL).Scan(
		&user.ID,
		&user.Username,
		&user.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) FindByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, avatar_url
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	err := ur.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.AvatarURL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) FindByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, avatar_url
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := ur.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.AvatarURL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
