package models

import "errors"

type Character struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Race         string `json:"race"`
	Class        string `json:"class"`
	Level        int    `json:"level"`
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Intelligence int    `json:"intelligence"`
	Wisdom       int    `json:"wisdom"`
	Charisma     int    `json:"charisma"`
	Photo        string `json:"photo"`
}

type CharacterShort struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type CharacterRepository interface {
	Create(character *Character) (*Character, error)
	GetAll() ([]CharacterShort, error)
	FindByID(id string) (*Character, error)
	Update(id string, update *Character) (*Character, error)
	Delete(id string) error
}

type CharacterService interface {
	Create(create *Character) (*Character, error)
	GetAll() ([]CharacterShort, error)
	FindByID(id string) (*Character, error)
	Update(id string, update *Character) (*Character, error)
	Delete(id string) error
}

var (
	ErrCharacterNotFound = errors.New("character not found")
)
