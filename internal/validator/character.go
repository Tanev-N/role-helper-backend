package validator

import (
	"fmt"
	"role-helper/internal/models"
	"strings"
)

const (
	MinAbilityScore = 1
	MaxAbilityScore = 30
	MinLevel        = 1
	MaxLevel        = 20
)


func ValidateCharacter(c models.Character) error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("имя не может быть пустым")
	}

	if strings.TrimSpace(c.Race) == "" {
		return fmt.Errorf("раса не может быть пустой")
	}

	if strings.TrimSpace(c.Class) == "" {
		return fmt.Errorf("класс не может быть пустым")
	}

	if c.Level < MinLevel || c.Level > MaxLevel {
		return fmt.Errorf("неверный уровень: %d. Уровень должен быть от %d до %d", c.Level, MinLevel, MaxLevel)
	}

	if c.Experience < 0 {
		return fmt.Errorf("опыт не может быть отрицательным")
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

	if c.ArmorClass < 0 {
		return fmt.Errorf("класс брони не может быть отрицательным")
	}
	if c.Speed < 0 {
		return fmt.Errorf("скорость не может быть отрицательной")
	}
	if c.HitPoints < 0 {
		return fmt.Errorf("хиты не могут быть отрицательными")
	}
	if c.MaxHitPoints < 0 {
		return fmt.Errorf("максимальные хиты не могут быть отрицательными")
	}
	if c.TempHitPoints < 0 {
		return fmt.Errorf("временные хиты не могут быть отрицательными")
	}

	return nil
}

func validateAbilityScore(name string, score int) error {
	if score < MinAbilityScore || score > MaxAbilityScore {
		return fmt.Errorf("неверная характеристика: %s равна %d. Характеристики должны быть от %d до %d",
			name, score, MinAbilityScore, MaxAbilityScore)
	}
	return nil
}
