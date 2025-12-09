package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"role-helper/internal/models"
)

type ArmorRepository struct {
	db *sql.DB
}

func NewArmorRepository(db *sql.DB) models.ArmorRepository {
	return &ArmorRepository{db: db}
}

func (ar *ArmorRepository) Create(armor *models.Armor) (*models.Armor, error) {
	query := `
		INSERT INTO armor (
			user_id, name, type, armor_class, modifier, cost, rarity,
			stealth_disadvantage, strength_requirement, weight, unique_stats, charges, modifiers, photo
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
		RETURNING id
	`

	err := ar.db.QueryRow(query,
		armor.UserID, armor.Name, armor.Type, armor.ArmorClass, armor.Modifier, armor.Cost, armor.Rarity,
		armor.StealthDisadvantage, armor.StrengthRequirement, armor.Weight, armor.UniqueStats, armor.Charges, armor.Modifiers, armor.Photo,
	).Scan(&armor.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create armor: %w", err)
	}

	return armor, nil
}

func (ar *ArmorRepository) GetAll(userID int) ([]models.Armor, error) {
	query := `
		SELECT id, user_id, name, type, armor_class, modifier, cost, rarity,
			stealth_disadvantage, strength_requirement, weight, unique_stats, charges, modifiers, photo
		FROM armor
		WHERE user_id = $1
		ORDER BY name
	`

	rows, err := ar.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query armor: %w", err)
	}
	defer rows.Close()

	var armors []models.Armor
	for rows.Next() {
		var armor models.Armor
		err := rows.Scan(
			&armor.ID, &armor.UserID, &armor.Name, &armor.Type, &armor.ArmorClass, &armor.Modifier,
			&armor.Cost, &armor.Rarity, &armor.StealthDisadvantage, &armor.StrengthRequirement,
			&armor.Weight, &armor.UniqueStats, &armor.Charges, &armor.Modifiers, &armor.Photo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan armor: %w", err)
		}
		armors = append(armors, armor)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating armor: %w", err)
	}

	return armors, nil
}

func (ar *ArmorRepository) GetByID(id, userID int) (*models.Armor, error) {
	query := `
		SELECT id, user_id, name, type, armor_class, modifier, cost, rarity,
			stealth_disadvantage, strength_requirement, weight, unique_stats, charges, modifiers, photo
		FROM armor
		WHERE id = $1 AND user_id = $2
	`

	var armor models.Armor
	err := ar.db.QueryRow(query, id, userID).Scan(
		&armor.ID, &armor.UserID, &armor.Name, &armor.Type, &armor.ArmorClass, &armor.Modifier,
		&armor.Cost, &armor.Rarity, &armor.StealthDisadvantage, &armor.StrengthRequirement,
		&armor.Weight, &armor.UniqueStats, &armor.Charges, &armor.Modifiers, &armor.Photo,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrArmorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get armor: %w", err)
	}

	return &armor, nil
}

func (ar *ArmorRepository) Update(id int, armor *models.Armor) (*models.Armor, error) {
	query := `
		UPDATE armor SET
			name = $2, type = $3, armor_class = $4, modifier = $5, cost = $6, rarity = $7,
			stealth_disadvantage = $8, strength_requirement = $9, weight = $10,
			unique_stats = $11, charges = $12, modifiers = $13, photo = $14
		WHERE id = $1
		RETURNING id, user_id, name, type, armor_class, modifier, cost, rarity,
			stealth_disadvantage, strength_requirement, weight, unique_stats, charges, modifiers, photo
	`

	updatedArmor := &models.Armor{}
	err := ar.db.QueryRow(query, id,
		armor.Name, armor.Type, armor.ArmorClass, armor.Modifier, armor.Cost, armor.Rarity,
		armor.StealthDisadvantage, armor.StrengthRequirement, armor.Weight,
		armor.UniqueStats, armor.Charges, armor.Modifiers, armor.Photo,
	).Scan(
		&updatedArmor.ID, &updatedArmor.UserID, &updatedArmor.Name, &updatedArmor.Type,
		&updatedArmor.ArmorClass, &updatedArmor.Modifier, &updatedArmor.Cost, &updatedArmor.Rarity,
		&updatedArmor.StealthDisadvantage, &updatedArmor.StrengthRequirement, &updatedArmor.Weight,
		&updatedArmor.UniqueStats, &updatedArmor.Charges, &updatedArmor.Modifiers, &updatedArmor.Photo,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrArmorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update armor: %w", err)
	}

	return updatedArmor, nil
}

func (ar *ArmorRepository) Delete(id, userID int) error {
	query := `DELETE FROM armor WHERE id = $1 AND user_id = $2`
	result, err := ar.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete armor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrArmorNotFound
	}

	return nil
}

func (ar *ArmorRepository) CheckBelonging(id, userID int) (bool, error) {
	query := `SELECT user_id FROM armor WHERE id = $1`
	var armorUserID int
	err := ar.db.QueryRow(query, id).Scan(&armorUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return armorUserID == userID, nil
}
