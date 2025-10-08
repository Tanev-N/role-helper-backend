package repository

import (
	"database/sql"
	"fmt"
	"role-helper/internal/models"
)

type CharacterDB struct {
	db *sql.DB
}

func NewCharacterDB(db *sql.DB) models.CharacterRepository {
	return &CharacterDB{db: db}
}

func (r *CharacterDB) Create(character *models.Character) (*models.Character, error) {
	query := `
		INSERT INTO characters (name, race, class, level, strength, dexterity, 
		                       constitution, intelligence, wisdom, charisma, photo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(query,
		character.Name,
		character.Race,
		character.Class,
		character.Level,
		character.Strength,
		character.Dexterity,
		character.Constitution,
		character.Intelligence,
		character.Wisdom,
		character.Charisma,
		character.Photo,
	).Scan(
		&id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create character: %w", err)
	}

	character.ID = id

	return character, nil
}

func (r *CharacterDB) GetAll() ([]models.CharacterShort, error) {
	query := `SELECT id, name, photo FROM characters`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get characters: %w", err)
	}
	defer rows.Close()

	var characters []models.CharacterShort
	for rows.Next() {
		var character models.CharacterShort
		err := rows.Scan(&character.ID, &character.Name, &character.Photo)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}
		characters = append(characters, character)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating characters: %w", err)
	}

	return characters, nil
}

func (r *CharacterDB) FindByID(id string) (*models.Character, error) {
	query := `
		SELECT id, name, race, class, level, strength, dexterity, 
		       constitution, intelligence, wisdom, charisma, photo
		FROM characters 
		WHERE id = $1
	`

	character := &models.Character{}
	err := r.db.QueryRow(query, id).Scan(
		&character.ID,
		&character.Name,
		&character.Race,
		&character.Class,
		&character.Level,
		&character.Strength,
		&character.Dexterity,
		&character.Constitution,
		&character.Intelligence,
		&character.Wisdom,
		&character.Charisma,
		&character.Photo,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed while searching character: %w", err)
	}

	return character, nil
}

func (r *CharacterDB) Update(id string, update *models.Character) (*models.Character, error) {
	query := `
		UPDATE characters 
		SET name = $1, race = $2, class = $3, level = $4, 
		    strength = $5, dexterity = $6, constitution = $7,
		    intelligence = $8, wisdom = $9, charisma = $10, photo = $11
		WHERE id = $12
	`

	_, err := r.db.Exec(query,
		update.Name,
		update.Race,
		update.Class,
		update.Level,
		update.Strength,
		update.Dexterity,
		update.Constitution,
		update.Intelligence,
		update.Wisdom,
		update.Charisma,
		update.Photo,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update character: %w", err)
	}

	return update, nil
}

func (r *CharacterDB) Delete(id string) error {
	query := `DELETE FROM characters WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	return nil
}
