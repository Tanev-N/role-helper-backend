package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"role-helper/internal/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo models.UserRepository
	redis    *redis.Client
}

type Session struct {
	UserID    int
	ExpiresAt time.Time
}

const sessionTTL = 24 * time.Hour

func NewUserUsecase(userRepo models.UserRepository, redisClient *redis.Client) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		redis:    redisClient,
	}
}

func (uu *UserUsecase) Register(req *models.UserRegisterRequest) (*models.User, string, error) {
	if req.Password != req.RePassword {
		return nil, "", models.ErrPasswordsDontMatch
	}
	if len(req.Username) < 4 {
		return nil, "", models.ErrInvalidCredentials
	}

	if len(req.Password) < 6 {
		return nil, "", models.ErrInvalidCredentials
	}

	_, err := uu.userRepo.FindByUsername(req.Username)
	if err == nil {
		return nil, "", models.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		AvatarURL:    "",
	}

	createdUser, err := uu.userRepo.Create(user)
	if err != nil {
		return nil, "", err
	}

	token, err := uu.generateToken()
	if err != nil {
		return nil, "", err
	}

	return createdUser, token, nil
}

func (uu *UserUsecase) Login(req *models.UserLoginRequest) (*models.User, string, error) {
	user, err := uu.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, "", models.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", models.ErrInvalidCredentials
	}

	token, err := uu.generateToken()
	if err != nil {
		return nil, "", err
	}

	ctx := context.Background()
	key := "session:" + token
	if err := uu.redis.Set(ctx, key, user.ID, sessionTTL).Err(); err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (uu *UserUsecase) Logout(token string) error {
	ctx := context.Background()
	key := "session:" + token
	return uu.redis.Del(ctx, key).Err()
}

func (uu *UserUsecase) ValidateToken(token string) (*models.User, error) {
	ctx := context.Background()
	key := "session:" + token
	val, err := uu.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, models.ErrInvalidToken
		}
		return nil, err
	}

	userID, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}

	user, err := uu.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uu *UserUsecase) generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
