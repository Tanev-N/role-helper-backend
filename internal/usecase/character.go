package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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

	if createReq.Skills == nil || len(createReq.Skills) == 0 {
		createReq.Skills = utils.GetDefaultSkills()
		utils.AutoCalculateCharacterStats(createReq)
	}

	// Если фото равно дефолтному значению, очищаем его
	if createReq.Photo == "/app/images/characters_default.png" || createReq.Photo == "characters_default.png" {
		createReq.Photo = ""
	}

	if err := validator.ValidateCharacter(*createReq); err != nil {
		return nil, err
	}
	return c.repo.Create(createReq)
}

func (c *CharacterUsecase) GetAll(userID int) ([]models.CharacterShort, error) {
	return c.repo.GetAll(userID)
}

func (c *CharacterUsecase) FindByID(id, userID int) (*models.Character, error) {
	// Проверяем, принадлежит ли персонаж пользователю
	belongs, err := c.repo.CheckBelonging(id, userID)
	if err != nil {
		return nil, err
	}

	if !belongs {
		inSameGame, err := c.repo.CheckCharacterInSameGame(id, userID)
		if err != nil {
			return nil, err
		}
		if !inSameGame {
			return nil, errors.New("character does not belong to user")
		}
	}

	character, err := c.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, models.ErrCharacterNotFound
	}
	return character, nil
}

func (c *CharacterUsecase) Update(id int, userID int, update *models.Character) (*models.Character, error) {
	ok, err := c.repo.CheckBelonging(id, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("character does not belong to user")
	}
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
		// Если фото равно дефолтному значению, очищаем его
		if update.Photo == "/app/images/characters_default.png" || update.Photo == "characters_default.png" {
			character.Photo = ""
		} else {
			character.Photo = update.Photo
		}
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

func (c *CharacterUsecase) Delete(id, userID int) error {
	ok, err := c.repo.CheckBelonging(id, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("character does not belong to user")
	}
	character, err := c.repo.FindByID(id)
	if err != nil {
		return err
	}
	if character == nil {
		return models.ErrCharacterNotFound
	}
	return c.repo.Delete(id)
}

func (c *CharacterUsecase) UploadPhoto(characterID, userID int, file multipart.File, originalFilename string) (string, error) {
	ok, err := c.repo.CheckBelonging(characterID, userID)
	if err != nil {
		return "", fmt.Errorf("ошибка проверки принадлежности: %w", err)
	}
	if !ok {
		return "", errors.New("персонаж не принадлежит пользователю")
	}

	character, err := c.repo.FindByID(characterID)
	if err != nil {
		return "", fmt.Errorf("персонаж не найден: %w", err)
	}
	if character == nil {
		return "", models.ErrCharacterNotFound
	}

	const (
		imageBaseURL = "https://critical-roll.ru/api/images"
		uploadDir    = "/var/www/app/images"
		maxFileSize  = 10 << 20 // 10 МБ
	)

	ext := strings.ToLower(filepath.Ext(originalFilename))
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !validExts[ext] {
		return "", fmt.Errorf("недопустимый формат изображения")
	}

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("ошибка генерации имени файла: %w", err)
	}
	filename := hex.EncodeToString(randomBytes) + ext
	filePath := filepath.Join(uploadDir, filename)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию: %w", err)
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer dst.Close()

	limitReader := io.LimitReader(file, maxFileSize+1)
	written, err := io.Copy(dst, limitReader)
	if err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("ошибка записи файла: %w", err)
	}
	if written > maxFileSize {
		os.Remove(filePath)
		return "", fmt.Errorf("файл превышает допустимый размер (макс. 10 МБ)")
	}

	photoURL := imageBaseURL + "/" + filename

	if err := c.repo.UpdatePhoto(character.ID, photoURL); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("ошибка обновления фото в БД: %w", err)
	}

	return photoURL, nil
}
