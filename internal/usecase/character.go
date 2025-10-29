package usecase

import (
	"fmt"
	"role-helper/internal/models"
	"role-helper/internal/utils"
	"role-helper/internal/validator"
	"strings"
)

type CharacterUsecase struct {
	repo models.CharacterRepository
}

func NewCharacterUsecase(repo models.CharacterRepository) models.CharacterService {
	return &CharacterUsecase{repo: repo}
}

func (c *CharacterUsecase) Create(createReq *models.Character) (*models.Character, error) {
	utils.AutoCalculateCharacterStats(createReq)
	
	if len(createReq.Skills) == 0 {
		createReq.Skills = utils.GetDefaultSkills()
		utils.AutoCalculateCharacterStats(createReq)
	}
	
	if err := validator.ValidateCharacter(*createReq); err != nil {
		return nil, err
	}
	return c.repo.Create(createReq)
}

func (c *CharacterUsecase) GetAll() ([]models.CharacterShort, error) {
	return c.repo.GetAll()
}

func (c *CharacterUsecase) FindByID(id string) (*models.Character, error) {
	character, err := c.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, models.ErrCharacterNotFound
	}
	return character, nil
}

func (c *CharacterUsecase) Update(id string, update *models.Character) (*models.Character, error) {
	character, err := c.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, models.ErrCharacterNotFound
	}

	if update.Name != "" && strings.TrimSpace(update.Name) == "" {
		return nil, fmt.Errorf("имя не может быть пустым")
	}
	if update.Race != "" && strings.TrimSpace(update.Race) == "" {
		return nil, fmt.Errorf("раса не может быть пустой")
	}
	if update.Class != "" && strings.TrimSpace(update.Class) == "" {
		return nil, fmt.Errorf("класс не может быть пустым")
	}
	if update.Level != 0 && (update.Level < 1 || update.Level > 20) {
		return nil, fmt.Errorf("неверный уровень: %d. Уровень должен быть от 1 до 20", update.Level)
	}
	if update.Experience < 0 {
		return nil, fmt.Errorf("опыт не может быть отрицательным")
	}

	if update.Strength != 0 && (update.Strength < 1 || update.Strength > 30) {
		return nil, fmt.Errorf("неверная характеристика: Сила равна %d. Характеристики должны быть от 1 до 30", update.Strength)
	}
	if update.Dexterity != 0 && (update.Dexterity < 1 || update.Dexterity > 30) {
		return nil, fmt.Errorf("неверная характеристика: Ловкость равна %d. Характеристики должны быть от 1 до 30", update.Dexterity)
	}
	if update.Constitution != 0 && (update.Constitution < 1 || update.Constitution > 30) {
		return nil, fmt.Errorf("неверная характеристика: Телосложение равна %d. Характеристики должны быть от 1 до 30", update.Constitution)
	}
	if update.Intelligence != 0 && (update.Intelligence < 1 || update.Intelligence > 30) {
		return nil, fmt.Errorf("неверная характеристика: Интеллект равна %d. Характеристики должны быть от 1 до 30", update.Intelligence)
	}
	if update.Wisdom != 0 && (update.Wisdom < 1 || update.Wisdom > 30) {
		return nil, fmt.Errorf("неверная характеристика: Мудрость равна %d. Характеристики должны быть от 1 до 30", update.Wisdom)
	}
	if update.Charisma != 0 && (update.Charisma < 1 || update.Charisma > 30) {
		return nil, fmt.Errorf("неверная характеристика: Харизма равна %d. Характеристики должны быть от 1 до 30", update.Charisma)
	}

	if update.Name != "" {
		character.Name = update.Name
	}
	if update.Race != "" {
		character.Race = update.Race
	}
	if update.Class != "" {
		character.Class = update.Class
	}
	if update.Level != 0 {
		character.Level = update.Level
	}
	if update.Alignment != "" {
		character.Alignment = update.Alignment
	}
	if update.Background != "" {
		character.Background = update.Background
	}
	if update.PlayerName != "" {
		character.PlayerName = update.PlayerName
	}
	if update.Experience != 0 {
		character.Experience = update.Experience
	}
	if update.Strength != 0 {
		character.Strength = update.Strength
	}
	if update.Dexterity != 0 {
		character.Dexterity = update.Dexterity
	}
	if update.Constitution != 0 {
		character.Constitution = update.Constitution
	}
	if update.Intelligence != 0 {
		character.Intelligence = update.Intelligence
	}
	if update.Wisdom != 0 {
		character.Wisdom = update.Wisdom
	}
	if update.Charisma != 0 {
		character.Charisma = update.Charisma
	}

	if update.ProficiencyBonus != 0 {
		character.ProficiencyBonus = update.ProficiencyBonus
	}
	if update.Initiative != 0 {
		character.Initiative = update.Initiative
	}
	if update.ArmorClass != 0 {
		character.ArmorClass = update.ArmorClass
	}
	if update.Speed != 0 {
		character.Speed = update.Speed
	}
	if update.HitPoints != 0 {
		character.HitPoints = update.HitPoints
	}
	if update.MaxHitPoints != 0 {
		character.MaxHitPoints = update.MaxHitPoints
	}
	if update.TempHitPoints != 0 {
		character.TempHitPoints = update.TempHitPoints
	}
	if update.HitDice != "" {
		character.HitDice = update.HitDice
	}
	if update.PersonalityTraits != "" {
		character.PersonalityTraits = update.PersonalityTraits
	}
	if update.Ideals != "" {
		character.Ideals = update.Ideals
	}
	if update.Bonds != "" {
		character.Bonds = update.Bonds
	}
	if update.Flaws != "" {
		character.Flaws = update.Flaws
	}
	if update.Proficiencies != "" {
		character.Proficiencies = update.Proficiencies
	}
	if update.Languages != "" {
		character.Languages = update.Languages
	}
	if update.Senses != "" {
		character.Senses = update.Senses
	}
	if update.Features != "" {
		character.Features = update.Features
	}
	if update.Photo != "" {
		character.Photo = update.Photo
	}
	if update.Skills != nil {
		character.Skills = update.Skills
	}
	if update.Equipment != nil {
		character.Equipment = update.Equipment
	}
	if update.Spells != nil {
		character.Spells = update.Spells
	}
	
	character.StrengthSave = update.StrengthSave
	character.DexteritySave = update.DexteritySave
	character.ConstitutionSave = update.ConstitutionSave
	character.IntelligenceSave = update.IntelligenceSave
	character.WisdomSave = update.WisdomSave
	character.CharismaSave = update.CharismaSave

	utils.AutoCalculateCharacterStats(character)
	
	return c.repo.Update(id, character)
}

func (c *CharacterUsecase) Delete(id string) error {
	character, err := c.repo.FindByID(id)
	if err != nil {
		return err
	}
	if character == nil {
		return models.ErrCharacterNotFound
	}
	return c.repo.Delete(id)
}
