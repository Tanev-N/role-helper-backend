package validator

import (
	"role-helper/internal/models"
	"strings"
)

func ValidateUserRegister(req *models.UserRegisterRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return models.ErrInvalidCredentials
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return models.ErrInvalidCredentials
	}

	if len(req.Password) < 6 {
		return models.ErrInvalidCredentials
	}

	if req.Password != req.RePassword {
		return models.ErrPasswordsDontMatch
	}

	return nil
}

func ValidateUserLogin(req *models.UserLoginRequest) error {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return models.ErrInvalidCredentials
	}

	return nil
}
