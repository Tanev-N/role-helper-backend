package utils

import (
	"math"
	"role-helper/internal/models"
	"strconv"
	"strings"
)


func CalculateAbilityModifier(abilityScore int) int {
	return int(math.Floor(float64(abilityScore-10) / 2))
}

func CalculateProficiencyBonus(level int) int {
	return 2 + (level-1)/4
}


func CalculateInitiative(dexterityMod int) int {
	return dexterityMod
}


func CalculateArmorClass(dexterityMod int, armorBonus int) int {
	return 10 + dexterityMod + armorBonus
}


func CalculateHitPoints(level int, hitDie int, constitutionMod int) int {	
	firstLevelHP := hitDie + constitutionMod
	
	averageHitDie := (hitDie + 1) / 2
	if (hitDie+1)%2 != 0 {
		averageHitDie++ 
	}
	additionalHP := averageHitDie + constitutionMod
	
	return firstLevelHP + (additionalHP * (level - 1))	
}

func CalculateSkillModifier(abilityMod int, proficient bool, proficiencyBonus int) int {
	if proficient {
		return abilityMod + proficiencyBonus
	}
	return abilityMod
}

func CalculateSavingThrowModifier(abilityMod int, proficient bool, proficiencyBonus int) int {
	if proficient {
		return abilityMod + proficiencyBonus
	}
	return abilityMod
}

func AutoCalculateCharacterStats(character *models.Character) {
	character.StrengthMod = CalculateAbilityModifier(character.Strength)
	character.DexterityMod = CalculateAbilityModifier(character.Dexterity)
	character.ConstitutionMod = CalculateAbilityModifier(character.Constitution)
	character.IntelligenceMod = CalculateAbilityModifier(character.Intelligence)
	character.WisdomMod = CalculateAbilityModifier(character.Wisdom)
	character.CharismaMod = CalculateAbilityModifier(character.Charisma)
	
	if character.ProficiencyBonus == 0 {
		character.ProficiencyBonus = CalculateProficiencyBonus(character.Level)
	}
	
	if character.Initiative == 0 {
		character.Initiative = CalculateInitiative(character.DexterityMod)
	}
	
	if character.MaxHitPoints == 0 {
		hitDie := parseHitDice(character.HitDice)
		character.MaxHitPoints = CalculateHitPoints(character.Level, hitDie, character.ConstitutionMod)
		character.HitPoints = character.MaxHitPoints
	}
	
	for i := range character.Skills {
		if character.Skills[i].Modifier == 0 {
			abilityMod := getAbilityModifierForSkill(character.Skills[i].Ability, character)
			character.Skills[i].Modifier = CalculateSkillModifier(abilityMod, character.Skills[i].Proficient, character.ProficiencyBonus)
		}
	}
}


func parseHitDice(hitDice string) int {
	if hitDice == "" || !strings.Contains(hitDice, "d") {
		return 8 
	}
	
	parts := strings.Split(hitDice, "d")
	if len(parts) != 2 {
		return 8
	}
	
	if dieSize, err := strconv.Atoi(parts[1]); err == nil {
		return dieSize
	}
	
	return 8 
}

func getAbilityModifierForSkill(ability string, character *models.Character) int {
	switch ability {
	case "Сила":
		return character.StrengthMod
	case "Ловкость":
		return character.DexterityMod
	case "Телосложение":
		return character.ConstitutionMod
	case "Интеллект":
		return character.IntelligenceMod
	case "Мудрость":
		return character.WisdomMod
	case "Харизма":
		return character.CharismaMod
	default:
		return 0
	}
}

func GetDefaultSkills() []models.CharacterSkill {
	return []models.CharacterSkill{
		{Name: "Атлетика", Ability: "Сила"},
		{Name: "Акробатика", Ability: "Ловкость"},
		{Name: "Ловкость рук", Ability: "Ловкость"},
		{Name: "Скрытность", Ability: "Ловкость"},
		{Name: "Анализ", Ability: "Интеллект"},
		{Name: "История", Ability: "Интеллект"},
		{Name: "Магия", Ability: "Интеллект"},
		{Name: "Природа", Ability: "Интеллект"},
		{Name: "Религия", Ability: "Интеллект"},
		{Name: "Внимательность", Ability: "Мудрость"},
		{Name: "Выживание", Ability: "Мудрость"},
		{Name: "Медицина", Ability: "Мудрость"},
		{Name: "Проницательность", Ability: "Мудрость"},
		{Name: "Уход за животными", Ability: "Мудрость"},
		{Name: "Выступление", Ability: "Харизма"},
		{Name: "Запугивание", Ability: "Харизма"},
		{Name: "Обман", Ability: "Харизма"},
		{Name: "Убеждение", Ability: "Харизма"},
	}
}
