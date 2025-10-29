package repository

import (
	"database/sql"
	"role-helper/internal/models"

	"github.com/google/uuid"
)

type CharacterRepository struct {
	db *sql.DB
}

func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

func (cr *CharacterRepository) Create(character *models.Character) (*models.Character, error) {
	character.ID = uuid.New().String()
	
	query := `
		INSERT INTO characters (
			id, name, race, class, level, alignment, background, player_name, experience,
			strength, dexterity, constitution, intelligence, wisdom, charisma,
			strength_mod, dexterity_mod, constitution_mod, intelligence_mod, wisdom_mod, charisma_mod,
			proficiency_bonus, initiative, armor_class, speed, hit_points, max_hit_points, temp_hit_points, hit_dice,
			strength_save, dexterity_save, constitution_save, intelligence_save, wisdom_save, charisma_save,
			personality_traits, ideals, bonds, flaws, proficiencies, languages, senses, features, photo
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21,
			$22, $23, $24, $25, $26, $27, $28, $29,
			$30, $31, $32, $33, $34, $35,
			$36, $37, $38, $39, $40, $41, $42, $43, $44
		)
	`

	_, err := cr.db.Exec(query,
		character.ID, character.Name, character.Race, character.Class, character.Level, character.Alignment, character.Background, character.PlayerName, character.Experience,
		character.Strength, character.Dexterity, character.Constitution, character.Intelligence, character.Wisdom, character.Charisma,
		character.StrengthMod, character.DexterityMod, character.ConstitutionMod, character.IntelligenceMod, character.WisdomMod, character.CharismaMod,
		character.ProficiencyBonus, character.Initiative, character.ArmorClass, character.Speed, character.HitPoints, character.MaxHitPoints, character.TempHitPoints, character.HitDice,
		character.StrengthSave, character.DexteritySave, character.ConstitutionSave, character.IntelligenceSave, character.WisdomSave, character.CharismaSave,
		character.PersonalityTraits, character.Ideals, character.Bonds, character.Flaws, character.Proficiencies, character.Languages, character.Senses, character.Features, character.Photo,
	)

	if err != nil {
		return nil, err
	}

	if err := cr.saveSkills(character.ID, character.Skills); err != nil {
		return nil, err
	}
	if err := cr.saveEquipment(character.ID, character.Equipment); err != nil {
		return nil, err
	}
	if err := cr.saveSpells(character.ID, character.Spells); err != nil {
		return nil, err
	}

	return character, nil
}

func (cr *CharacterRepository) GetAll() ([]models.CharacterShort, error) {
	query := `SELECT id, name, photo FROM characters ORDER BY name`
	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []models.CharacterShort
	for rows.Next() {
		var character models.CharacterShort
		err := rows.Scan(&character.ID, &character.Name, &character.Photo)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}

func (cr *CharacterRepository) FindByID(id string) (*models.Character, error) {
	query := `
		SELECT id, name, race, class, level, alignment, background, player_name, experience,
			strength, dexterity, constitution, intelligence, wisdom, charisma,
			strength_mod, dexterity_mod, constitution_mod, intelligence_mod, wisdom_mod, charisma_mod,
			proficiency_bonus, initiative, armor_class, speed, hit_points, max_hit_points, temp_hit_points, hit_dice,
			strength_save, dexterity_save, constitution_save, intelligence_save, wisdom_save, charisma_save,
			personality_traits, ideals, bonds, flaws, proficiencies, languages, senses, features, photo
		FROM characters WHERE id = $1
	`

	character := &models.Character{}
	err := cr.db.QueryRow(query, id).Scan(
		&character.ID, &character.Name, &character.Race, &character.Class, &character.Level, &character.Alignment, &character.Background, &character.PlayerName, &character.Experience,
		&character.Strength, &character.Dexterity, &character.Constitution, &character.Intelligence, &character.Wisdom, &character.Charisma,
		&character.StrengthMod, &character.DexterityMod, &character.ConstitutionMod, &character.IntelligenceMod, &character.WisdomMod, &character.CharismaMod,
		&character.ProficiencyBonus, &character.Initiative, &character.ArmorClass, &character.Speed, &character.HitPoints, &character.MaxHitPoints, &character.TempHitPoints, &character.HitDice,
		&character.StrengthSave, &character.DexteritySave, &character.ConstitutionSave, &character.IntelligenceSave, &character.WisdomSave, &character.CharismaSave,
		&character.PersonalityTraits, &character.Ideals, &character.Bonds, &character.Flaws, &character.Proficiencies, &character.Languages, &character.Senses, &character.Features, &character.Photo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrCharacterNotFound
		}
		return nil, err
	}

	character.Skills, err = cr.getSkills(character.ID)
	if err != nil {
		return nil, err
	}

	character.Equipment, err = cr.getEquipment(character.ID)
	if err != nil {
		return nil, err
	}

	character.Spells, err = cr.getSpells(character.ID)
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (cr *CharacterRepository) Update(id string, update *models.Character) (*models.Character, error) {
	existing, err := cr.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, models.ErrCharacterNotFound
	}

	// Обновляем только непустые поля
	if update.Name != "" {
		existing.Name = update.Name
	}
	if update.Race != "" {
		existing.Race = update.Race
	}
	if update.Class != "" {
		existing.Class = update.Class
	}
	if update.Level != 0 {
		existing.Level = update.Level
	}
	if update.Alignment != "" {
		existing.Alignment = update.Alignment
	}
	if update.Background != "" {
		existing.Background = update.Background
	}
	if update.PlayerName != "" {
		existing.PlayerName = update.PlayerName
	}
	if update.Experience != 0 {
		existing.Experience = update.Experience
	}

	if update.Strength != 0 {
		existing.Strength = update.Strength
	}
	if update.Dexterity != 0 {
		existing.Dexterity = update.Dexterity
	}
	if update.Constitution != 0 {
		existing.Constitution = update.Constitution
	}
	if update.Intelligence != 0 {
		existing.Intelligence = update.Intelligence
	}
	if update.Wisdom != 0 {
		existing.Wisdom = update.Wisdom
	}
	if update.Charisma != 0 {
		existing.Charisma = update.Charisma
	}

	existing.StrengthMod = update.StrengthMod
	existing.DexterityMod = update.DexterityMod
	existing.ConstitutionMod = update.ConstitutionMod
	existing.IntelligenceMod = update.IntelligenceMod
	existing.WisdomMod = update.WisdomMod
	existing.CharismaMod = update.CharismaMod

	if update.ProficiencyBonus != 0 {
		existing.ProficiencyBonus = update.ProficiencyBonus
	}
	if update.Initiative != 0 {
		existing.Initiative = update.Initiative
	}
	if update.ArmorClass != 0 {
		existing.ArmorClass = update.ArmorClass
	}
	if update.Speed != 0 {
		existing.Speed = update.Speed
	}
	if update.HitPoints != 0 {
		existing.HitPoints = update.HitPoints
	}
	if update.MaxHitPoints != 0 {
		existing.MaxHitPoints = update.MaxHitPoints
	}
	if update.TempHitPoints != 0 {
		existing.TempHitPoints = update.TempHitPoints
	}
	if update.HitDice != "" {
		existing.HitDice = update.HitDice
	}

	existing.StrengthSave = update.StrengthSave
	existing.DexteritySave = update.DexteritySave
	existing.ConstitutionSave = update.ConstitutionSave
	existing.IntelligenceSave = update.IntelligenceSave
	existing.WisdomSave = update.WisdomSave
	existing.CharismaSave = update.CharismaSave

	if update.PersonalityTraits != "" {
		existing.PersonalityTraits = update.PersonalityTraits
	}
	if update.Ideals != "" {
		existing.Ideals = update.Ideals
	}
	if update.Bonds != "" {
		existing.Bonds = update.Bonds
	}
	if update.Flaws != "" {
		existing.Flaws = update.Flaws
	}
	if update.Proficiencies != "" {
		existing.Proficiencies = update.Proficiencies
	}
	if update.Languages != "" {
		existing.Languages = update.Languages
	}
	if update.Senses != "" {
		existing.Senses = update.Senses
	}
	if update.Features != "" {
		existing.Features = update.Features
	}
	if update.Photo != "" {
		existing.Photo = update.Photo
	}

	if update.Skills != nil {
		existing.Skills = update.Skills
	}
	if update.Equipment != nil {
		existing.Equipment = update.Equipment
	}
	if update.Spells != nil {
		existing.Spells = update.Spells
	}

	query := `
		UPDATE characters SET
			name = $2, race = $3, class = $4, level = $5, alignment = $6, background = $7, player_name = $8, experience = $9,
			strength = $10, dexterity = $11, constitution = $12, intelligence = $13, wisdom = $14, charisma = $15,
			strength_mod = $16, dexterity_mod = $17, constitution_mod = $18, intelligence_mod = $19, wisdom_mod = $20, charisma_mod = $21,
			proficiency_bonus = $22, initiative = $23, armor_class = $24, speed = $25, hit_points = $26, max_hit_points = $27, temp_hit_points = $28, hit_dice = $29,
			strength_save = $30, dexterity_save = $31, constitution_save = $32, intelligence_save = $33, wisdom_save = $34, charisma_save = $35,
			personality_traits = $36, ideals = $37, bonds = $38, flaws = $39, proficiencies = $40, languages = $41, senses = $42, features = $43, photo = $44
		WHERE id = $1
	`

	_, err = cr.db.Exec(query, id,
		existing.Name, existing.Race, existing.Class, existing.Level, existing.Alignment, existing.Background, existing.PlayerName, existing.Experience,
		existing.Strength, existing.Dexterity, existing.Constitution, existing.Intelligence, existing.Wisdom, existing.Charisma,
		existing.StrengthMod, existing.DexterityMod, existing.ConstitutionMod, existing.IntelligenceMod, existing.WisdomMod, existing.CharismaMod,
		existing.ProficiencyBonus, existing.Initiative, existing.ArmorClass, existing.Speed, existing.HitPoints, existing.MaxHitPoints, existing.TempHitPoints, existing.HitDice,
		existing.StrengthSave, existing.DexteritySave, existing.ConstitutionSave, existing.IntelligenceSave, existing.WisdomSave, existing.CharismaSave,
		existing.PersonalityTraits, existing.Ideals, existing.Bonds, existing.Flaws, existing.Proficiencies, existing.Languages, existing.Senses, existing.Features, existing.Photo,
	)

	if err != nil {
		return nil, err
	}

	if update.Skills != nil {
		if err := cr.updateSkills(existing.ID, existing.Skills); err != nil {
			return nil, err
		}
	}
	if update.Equipment != nil {
		if err := cr.updateEquipment(existing.ID, existing.Equipment); err != nil {
			return nil, err
		}
	}
	if update.Spells != nil {
		if err := cr.updateSpells(existing.ID, existing.Spells); err != nil {
			return nil, err
		}
	}

	return existing, nil
}

func (cr *CharacterRepository) Delete(id string) error {
	query := `DELETE FROM characters WHERE id = $1`
	result, err := cr.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrCharacterNotFound
	}

	return nil
}


func (cr *CharacterRepository) getSkills(characterID string) ([]models.CharacterSkill, error) {
	query := `SELECT name, modifier, proficient, ability FROM character_skills WHERE character_id = $1 ORDER BY name`
	rows, err := cr.db.Query(query, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.CharacterSkill
	for rows.Next() {
		var skill models.CharacterSkill
		err := rows.Scan(&skill.Name, &skill.Modifier, &skill.Proficient, &skill.Ability)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (cr *CharacterRepository) saveSkills(characterID string, skills []models.CharacterSkill) error {
	if len(skills) == 0 {
		return nil
	}

	query := `INSERT INTO character_skills (character_id, name, modifier, proficient, ability) VALUES ($1, $2, $3, $4, $5)`
	for _, skill := range skills {
		_, err := cr.db.Exec(query, characterID, skill.Name, skill.Modifier, skill.Proficient, skill.Ability)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cr *CharacterRepository) updateSkills(characterID string, skills []models.CharacterSkill) error {
	deleteQuery := `DELETE FROM character_skills WHERE character_id = $1`
	_, err := cr.db.Exec(deleteQuery, characterID)
	if err != nil {
		return err
	}

	return cr.saveSkills(characterID, skills)
}

func (cr *CharacterRepository) getEquipment(characterID string) ([]models.Equipment, error) {
	query := `SELECT name, description FROM character_equipment WHERE character_id = $1 ORDER BY name`
	rows, err := cr.db.Query(query, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipment []models.Equipment
	for rows.Next() {
		var item models.Equipment
		err := rows.Scan(&item.Name, &item.Description)
		if err != nil {
			return nil, err
		}
		equipment = append(equipment, item)
	}

	return equipment, nil
}

func (cr *CharacterRepository) saveEquipment(characterID string, equipment []models.Equipment) error {
	if len(equipment) == 0 {
		return nil
	}

	query := `INSERT INTO character_equipment (character_id, name, description) VALUES ($1, $2, $3)`
	for _, item := range equipment {
		_, err := cr.db.Exec(query, characterID, item.Name, item.Description)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cr *CharacterRepository) updateEquipment(characterID string, equipment []models.Equipment) error {
	deleteQuery := `DELETE FROM character_equipment WHERE character_id = $1`
	_, err := cr.db.Exec(deleteQuery, characterID)
	if err != nil {
		return err
	}

	return cr.saveEquipment(characterID, equipment)
}

func (cr *CharacterRepository) getSpells(characterID string) ([]models.Spell, error) {
	query := `SELECT name, description FROM character_spells WHERE character_id = $1 ORDER BY name`
	rows, err := cr.db.Query(query, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []models.Spell
	for rows.Next() {
		var spell models.Spell
		err := rows.Scan(&spell.Name, &spell.Description)
		if err != nil {
			return nil, err
		}
		spells = append(spells, spell)
	}

	return spells, nil
}

func (cr *CharacterRepository) saveSpells(characterID string, spells []models.Spell) error {
	if len(spells) == 0 {
		return nil
	}

	query := `INSERT INTO character_spells (character_id, name, description) VALUES ($1, $2, $3)`
	for _, spell := range spells {
		_, err := cr.db.Exec(query, characterID, spell.Name, spell.Description)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cr *CharacterRepository) updateSpells(characterID string, spells []models.Spell) error {
	deleteQuery := `DELETE FROM character_spells WHERE character_id = $1`
	_, err := cr.db.Exec(deleteQuery, characterID)
	if err != nil {
		return err
	}
	return cr.saveSpells(characterID, spells)
}