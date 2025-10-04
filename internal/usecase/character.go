package usecase

import (
	"role-helper/internal/models"
)

type CharacterUsecase struct {
	repo models.CharacterRepository
}

func NewCharacterUsecase(repo models.CharacterRepository) models.CharacterService {
	return &CharacterUsecase{repo: repo}
}

func (c *CharacterUsecase) Create(createReq *models.Character) (*models.Character, error) {
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
	return c.repo.Update(id, update)
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
