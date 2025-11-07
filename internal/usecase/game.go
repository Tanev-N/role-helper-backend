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

func NewGameUseCase(game *repository.GameRepository, user *repository.UserRepository, char models.CharacterRepository) models.GameService {
	return &GameUseCase{game: game, user: user, character: char}
}

func (uc *GameUseCase) CreateGame(gameReq *models.Game) (*models.Game, error) {
	return uc.game.CreateGame(gameReq)
}

func (uc *GameUseCase) CreateSession(session *models.Session) (*models.Session, []models.Session, error) {
	session, err := uc.game.CreateSession(session)
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

	if character.PlayerName != user.Username {
		return nil, errors.New("invalid user id for character")
	}

	session, err := uc.game.GetSessionByKey(key)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	player.GameID = session.GameID

	err = uc.game.AddPlayerToGame(player)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *GameUseCase) FinishSession(id int, summary string) error {
	session, err := uc.game.GetSessionByID(id)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("invalid session")
	}

	return uc.game.FinishSession(id, summary)
}

func (uc *GameUseCase) GetAllGames(userID int) ([]models.Game, error) {
	games, err := uc.game.GetAllGames(userID)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (uc *GameUseCase) GetGamePlayers(id int) ([]models.GamePlayer, error) {
	players, err := uc.game.GetGamePlayers(id)
	if err != nil {
		return nil, err
	}
	return players, nil
}
