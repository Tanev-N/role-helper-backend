package repository

import (
	"github.com/google/uuid"
	"role-helper/internal/models"
)

type CharacterMemory struct {
	characters map[string]*models.Character
}

func NewCharacterMemory() models.CharacterRepository {
	return &CharacterMemory{
		characters: make(map[string]*models.Character),
	}
}

func (r *CharacterMemory) Create(character *models.Character) (*models.Character, error) {
	character.ID = uuid.New().String()
	r.characters[character.ID] = character
	return character, nil
}

func (r *CharacterMemory) GetAll() ([]models.CharacterShort, error) {
	characters := make([]models.CharacterShort, 0, len(r.characters))
	for _, character := range r.characters {
		characters = append(characters, models.CharacterShort{ID: character.ID, Name: character.Name})
	}
	return characters, nil
}

func (r *CharacterMemory) FindByID(id string) (*models.Character, error) {
	character, exists := r.characters[id]
	if !exists {
		return nil, nil
	}
	return character, nil
}

func (r *CharacterMemory) Update(id string, update *models.Character) (*models.Character, error) {
	character := r.characters[id]

	character.Name = update.Name
	character.Race = update.Race
	character.Class = update.Class
	character.Level = update.Level
	character.Strength = update.Strength
	character.Dexterity = update.Dexterity
	character.Constitution = update.Constitution
	character.Intelligence = update.Intelligence
	character.Wisdom = update.Wisdom
	character.Charisma = update.Charisma
	character.Photo = update.Photo

	return character, nil
}

func (r *CharacterMemory) Delete(id string) error {
	delete(r.characters, id)
	return nil
}
