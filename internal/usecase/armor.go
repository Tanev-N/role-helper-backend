package usecase

import (
	"errors"
	"role-helper/internal/models"
)

type ArmorUsecase struct {
	repo models.ArmorRepository
}

func NewArmorUsecase(repo models.ArmorRepository) models.ArmorService {
	return &ArmorUsecase{repo: repo}
}

func (a *ArmorUsecase) Create(armor *models.Armor) (*models.Armor, error) {
	if armor.Name == "" {
		return nil, errors.New("name is required")
	}
	return a.repo.Create(armor)
}

func (a *ArmorUsecase) GetAll(userID int) ([]models.Armor, error) {
	return a.repo.GetAll(userID)
}

func (a *ArmorUsecase) GetByID(id, userID int) (*models.Armor, error) {
	ok, err := a.repo.CheckBelonging(id, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("armor does not belong to user")
	}
	return a.repo.GetByID(id, userID)
}

func (a *ArmorUsecase) Update(id int, userID int, armor *models.Armor) (*models.Armor, error) {
	ok, err := a.repo.CheckBelonging(id, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("armor does not belong to user")
	}
	return a.repo.Update(id, armor)
}

func (a *ArmorUsecase) Delete(id, userID int) error {
	ok, err := a.repo.CheckBelonging(id, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("armor does not belong to user")
	}
	return a.repo.Delete(id, userID)
}
