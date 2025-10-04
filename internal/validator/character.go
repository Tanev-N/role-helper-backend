package validator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrCharacterNotFound = errors.New("Персонаж не найден")
)

const (
	MinAbilityScore = 1
	MaxAbilityScore = 30
	MinLevel        = 1
	MaxLevel        = 20
)

var ValidRaces = []string{
	"Человек", "Эльф", "Дварф", "Халфлинг", "Гном", "Полуорк", "Полуэльф", "Тифлинг", "Драконорожденный",
}

var ValidClasses = []string{
	"Воин", "Маг", "Плут", "Жрец", "Следопыт", "Паладин", "Варвар", "Бард", "Друид", "Монах", "Чародей", "Колдун",
}

type Character struct {
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

func ValidateCharacter(c Character) error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("имя не может быть пустым")
	}

	if !isValidRace(c.Race) {
		return fmt.Errorf("неверная раса: %s. Доступные расы: %s", c.Race, strings.Join(ValidRaces, ", "))
	}

	if !isValidClass(c.Class) {
		return fmt.Errorf("неверный класс: %s. Доступные классы: %s", c.Class, strings.Join(ValidClasses, ", "))
	}

	if c.Level < MinLevel || c.Level > MaxLevel {
		return fmt.Errorf("неверный уровень: %d. Уровень должен быть от %d до %d", c.Level, MinLevel, MaxLevel)
	}

	if err := validateAbilityScore("Сила", c.Strength); err != nil {
		return err
	}
	if err := validateAbilityScore("Ловкость", c.Dexterity); err != nil {
		return err
	}
	if err := validateAbilityScore("Телосложение", c.Constitution); err != nil {
		return err
	}
	if err := validateAbilityScore("Интеллект", c.Intelligence); err != nil {
		return err
	}
	if err := validateAbilityScore("Мудрость", c.Wisdom); err != nil {
		return err
	}
	if err := validateAbilityScore("Харизма", c.Charisma); err != nil {
		return err
	}

	return nil
}

func isValidRace(race string) bool {
	for _, validRace := range ValidRaces {
		if strings.EqualFold(race, validRace) {
			return true
		}
	}
	return false
}

func isValidClass(class string) bool {
	for _, validClass := range ValidClasses {
		if strings.EqualFold(class, validClass) {
			return true
		}
	}
	return false
}

func validateAbilityScore(name string, score int) error {
	if score < MinAbilityScore || score > MaxAbilityScore {
		return fmt.Errorf("неверная характеристика: %s равна %d. Характеристики должны быть от %d до %d",
			name, score, MinAbilityScore, MaxAbilityScore)
	}
	return nil
}
