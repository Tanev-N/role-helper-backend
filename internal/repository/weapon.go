package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"role-helper/internal/models"
)

type WeaponRepository struct {
	db *sql.DB
}

func NewWeaponRepository(db *sql.DB) models.WeaponRepository {
	return &WeaponRepository{db: db}
}

func (wr *WeaponRepository) Create(weapon *models.Weapon) (*models.Weapon, error) {
	query := `
		INSERT INTO weapons (
			user_id, name, type, damage, modifier, cost, rarity,
			grip, range_meters, weight, unique_stats, charges, photo
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
		RETURNING id
	`

	err := wr.db.QueryRow(query,
		weapon.UserID, weapon.Name, weapon.Type, weapon.Damage, weapon.Modifier, weapon.Cost, weapon.Rarity,
		weapon.Grip, weapon.RangeMeters, weapon.Weight, weapon.UniqueStats, weapon.Charges, weapon.Photo,
	).Scan(&weapon.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create weapon: %w", err)
	}

	return weapon, nil
}

func (wr *WeaponRepository) GetAll(userID int) ([]models.Weapon, error) {
	query := `
		SELECT id, user_id, name, type, damage, modifier, cost, rarity,
			grip, range_meters, weight, unique_stats, charges, photo
		FROM weapons
		WHERE user_id = $1
		ORDER BY name
	`

	rows, err := wr.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query weapons: %w", err)
	}
	defer rows.Close()

	var weapons []models.Weapon
	for rows.Next() {
		var weapon models.Weapon
		err := rows.Scan(
			&weapon.ID, &weapon.UserID, &weapon.Name, &weapon.Type, &weapon.Damage, &weapon.Modifier,
			&weapon.Cost, &weapon.Rarity, &weapon.Grip, &weapon.RangeMeters,
			&weapon.Weight, &weapon.UniqueStats, &weapon.Charges, &weapon.Photo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weapon: %w", err)
		}
		weapons = append(weapons, weapon)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating weapons: %w", err)
	}

	return weapons, nil
}

func (wr *WeaponRepository) GetByID(id, userID int) (*models.Weapon, error) {
	query := `
		SELECT id, user_id, name, type, damage, modifier, cost, rarity,
			grip, range_meters, weight, unique_stats, charges, photo
		FROM weapons
		WHERE id = $1 AND user_id = $2
	`

	var weapon models.Weapon
	err := wr.db.QueryRow(query, id, userID).Scan(
		&weapon.ID, &weapon.UserID, &weapon.Name, &weapon.Type, &weapon.Damage, &weapon.Modifier,
		&weapon.Cost, &weapon.Rarity, &weapon.Grip, &weapon.RangeMeters,
		&weapon.Weight, &weapon.UniqueStats, &weapon.Charges, &weapon.Photo,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrWeaponNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get weapon: %w", err)
	}

	return &weapon, nil
}

func (wr *WeaponRepository) Update(id int, weapon *models.Weapon) (*models.Weapon, error) {
	query := `
		UPDATE weapons SET
			name = $2, type = $3, damage = $4, modifier = $5, cost = $6, rarity = $7,
			grip = $8, range_meters = $9, weight = $10,
			unique_stats = $11, charges = $12, photo = $13
		WHERE id = $1
		RETURNING id, user_id, name, type, damage, modifier, cost, rarity,
			grip, range_meters, weight, unique_stats, charges, photo
	`

	updatedWeapon := &models.Weapon{}
	err := wr.db.QueryRow(query, id,
		weapon.Name, weapon.Type, weapon.Damage, weapon.Modifier, weapon.Cost, weapon.Rarity,
		weapon.Grip, weapon.RangeMeters, weapon.Weight,
		weapon.UniqueStats, weapon.Charges, weapon.Photo,
	).Scan(
		&updatedWeapon.ID, &updatedWeapon.UserID, &updatedWeapon.Name, &updatedWeapon.Type,
		&updatedWeapon.Damage, &updatedWeapon.Modifier, &updatedWeapon.Cost, &updatedWeapon.Rarity,
		&updatedWeapon.Grip, &updatedWeapon.RangeMeters, &updatedWeapon.Weight,
		&updatedWeapon.UniqueStats, &updatedWeapon.Charges, &updatedWeapon.Photo,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrWeaponNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update weapon: %w", err)
	}

	return updatedWeapon, nil
}

func (wr *WeaponRepository) Delete(id, userID int) error {
	query := `DELETE FROM weapons WHERE id = $1 AND user_id = $2`
	result, err := wr.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete weapon: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrWeaponNotFound
	}

	return nil
}

func (wr *WeaponRepository) CheckBelonging(id, userID int) (bool, error) {
	query := `SELECT user_id FROM weapons WHERE id = $1`
	var weaponUserID int
	err := wr.db.QueryRow(query, id).Scan(&weaponUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return weaponUserID == userID, nil
}
