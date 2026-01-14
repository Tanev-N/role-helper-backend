package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"role-helper/internal/models"
	"strconv"
	"strings"
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

const (
	imageBaseURL = "https://critical-roll.ru/api/images"
	uploadDir    = "/var/www/app/images"
	maxFileSize  = 10 << 20 // 10 МБ
)

func (uu *UserUsecase) UploadAvatar(userID int, file multipart.File, originalFilename string) (string, error) {
	user, err := uu.userRepo.FindByID(userID)
	if err != nil {
		return "", fmt.Errorf("пользователь не найден: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(originalFilename))
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !validExts[ext] {
		return "", fmt.Errorf("недопустимый формат изображения")
	}

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("ошибка генерации имени файла: %w", err)
	}
	filename := hex.EncodeToString(randomBytes) + ext
	filePath := filepath.Join(uploadDir, filename)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию: %w", err)
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer dst.Close()

	limitReader := io.LimitReader(file, maxFileSize+1)
	written, err := io.Copy(dst, limitReader)
	if err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("ошибка записи файла: %w", err)
	}
	if written > maxFileSize {
		os.Remove(filePath)
		return "", fmt.Errorf("файл превышает допустимый размер (макс. 10 МБ)")
	}

	avatarURL := imageBaseURL + "/" + filename

	if err := uu.userRepo.UpdateAvatar(user.ID, avatarURL); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("ошибка обновления аватара в БД: %w", err)
	}

	return avatarURL, nil
}

func (uu *UserUsecase) Register(req *models.UserRegisterRequest) (*models.User, string, error) {
	_, err := uu.userRepo.FindByUsername(req.Username)
	if err == nil {
		return nil, "", models.ErrUserAlreadyExists
	}
	if err != nil && err != models.ErrUserNotFound {
		return nil, "", err
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

	ctx := context.Background()
	key := "session:" + token
	id := strconv.Itoa(createdUser.ID)
	if err := uu.redis.Set(ctx, key, id, sessionTTL).Err(); err != nil {
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
	id := strconv.Itoa(user.ID)
	if err := uu.redis.Set(ctx, key, id, sessionTTL).Err(); err != nil {
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
