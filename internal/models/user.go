package models

import (
	"errors"
	"mime/multipart"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	AvatarURL    string `json:"avatar_url"`
}

type UserRegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepository interface {
	Create(user *User) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByID(id int) (*User, error)
	UpdateAvatar(userID int, avatarURL string) error
}

type UserService interface {
	Register(req *UserRegisterRequest) (*User, string, error)
	Login(req *UserLoginRequest) (*User, string, error)
	Logout(token string) error
	ValidateToken(token string) (*User, error)
	UploadAvatar(userID int, file multipart.File, originalFilename string) (string, error)
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrPasswordsDontMatch = errors.New("passwords don't match")
	ErrInvalidToken       = errors.New("invalid token")
)
