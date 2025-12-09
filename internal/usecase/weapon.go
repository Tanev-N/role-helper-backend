package usecase

import (
	"errors"
	"role-helper/internal/models"
)

type WeaponUsecase struct {
	repo models.WeaponRepository
}

func NewWeaponUsecase(repo models.WeaponRepository) models.WeaponService {
	return &WeaponUsecase{repo: repo}
}

func (w *WeaponUsecase) Create(weapon *models.Weapon) (*models.Weapon, error) {
	if weapon.Name == "" {
		return nil, errors.New("name is required")
	}
	return w.repo.Create(weapon)
}

func (w *WeaponUsecase) GetAll(userID int) ([]models.Weapon, error) {
	return w.repo.GetAll(userID)
}

func (w *WeaponUsecase) GetByID(id, userID int) (*models.Weapon, error) {
	ok, err := w.repo.CheckBelonging(id, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("weapon does not belong to user")
	}
	return w.repo.GetByID(id, userID)
}

func (w *WeaponUsecase) Update(id int, userID int, weapon *models.Weapon) (*models.Weapon, error) {
	ok, err := w.repo.CheckBelonging(id, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("weapon does not belong to user")
	}
	return w.repo.Update(id, weapon)
}

func (w *WeaponUsecase) Delete(id, userID int) error {
	ok, err := w.repo.CheckBelonging(id, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("weapon does not belong to user")
	}
	return w.repo.Delete(id, userID)
}
