package usecase

import (
	"role-helper/internal/models"
	"role-helper/internal/validator"
)

type CharacterUsecase struct {
	repo models.CharacterRepository
}

func NewCharacterUsecase(repo models.CharacterRepository) models.CharacterService {
	return &CharacterUsecase{repo: repo}
}

func (c *CharacterUsecase) Create(createReq *models.Character) (*models.Character, error) {
	characterForValidation := validator.Character{
		Name:         createReq.Name,
		Race:         createReq.Race,
		Class:        createReq.Class,
		Level:        createReq.Level,
		Strength:     createReq.Strength,
		Dexterity:    createReq.Dexterity,
		Constitution: createReq.Constitution,
		Intelligence: createReq.Intelligence,
		Wisdom:       createReq.Wisdom,
		Charisma:     createReq.Charisma,
		Photo:        createReq.Photo,
	}

	if err := validator.ValidateCharacter(characterForValidation); err != nil {
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
		return nil, validator.ErrCharacterNotFound
	}
	return character, nil
}

func (c *CharacterUsecase) Update(id string, update *models.Character) (*models.Character, error) {
	character, err := c.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, validator.ErrCharacterNotFound
	}

	characterForValidation := validator.Character{
		Name:         update.Name,
		Race:         update.Race,
		Class:        update.Class,
		Level:        update.Level,
		Strength:     update.Strength,
		Dexterity:    update.Dexterity,
		Constitution: update.Constitution,
		Intelligence: update.Intelligence,
		Wisdom:       update.Wisdom,
		Charisma:     update.Charisma,
		Photo:        update.Photo,
	}

	if err := validator.ValidateCharacter(characterForValidation); err != nil {
		return nil, err
	}

	return c.repo.Update(id, update)
}

func (c *CharacterUsecase) Delete(id string) error {
	character, err := c.repo.FindByID(id)
	if err != nil {
		return err
	}
	if character == nil {
		return validator.ErrCharacterNotFound
	}
	return c.repo.Delete(id)
}
