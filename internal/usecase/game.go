package usecase

import (
	"errors"
	"role-helper/internal/models"
	"role-helper/internal/repository"
)

type GameUseCase struct {
	game      models.GameRepository
	user      models.UserRepository
	character models.CharacterRepository
}

func NewGameUseCase(game *repository.GameRepository, user *repository.UserRepository) models.GameService {
	return &GameUseCase{game: game, user: user}
}

func (uc *GameUseCase) CreateGame(gameReq *models.Game) (*models.Game, error) {
	if gameReq.Name == "" {
		return nil, errors.New("game name cannot be empty")
	}
	user, err := uc.user.FindByID(gameReq.MasterID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid master id")
	}
	return uc.game.CreateGame(gameReq)
}

func (uc *GameUseCase) CreateSession(session *models.Session) (*models.Session, []models.Session, error) {
	game, err := uc.game.GetGameByID(session.GameID)
	if err != nil {
		return nil, nil, err
	}
	if game == nil {
		return nil, nil, errors.New("game not found")
	}

	session, err = uc.game.CreateSession(session)
	if err != nil {
		return nil, nil, err
	}

	sessions, err := uc.game.GetPreviousSessions(session.GameID)
	if err != nil {
		return nil, nil, err
	}

	return session, sessions, nil
}

func (uc *GameUseCase) EnterSession(key string, player *models.GamePlayer) (*models.Session, error) {
	if key == "" {
		return nil, errors.New("empty session key")
	}

	user, err := uc.user.FindByID(player.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid player id")
	}

	character, err := uc.character.FindByID(player.CharacterID)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, errors.New("invalid character id")
	}

	session, err := uc.game.GetSessionByKey(key)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	err = uc.game.AddPlayerToGame(player)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *GameUseCase) FinishSession(id int, summary string) error {
	if summary == "" {
		return errors.New("summary cannot be empty")
	}

	return uc.game.FinishSession(id, summary)
}
