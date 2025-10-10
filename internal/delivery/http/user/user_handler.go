package user

import (
	"encoding/json"
	"log"
	"net/http"
	"role-helper/internal/models"
	"role-helper/internal/validator"
	"time"
)

type UserRouter struct {
	UserUsecase models.UserService
}

func NewUserRouter(userUsecase models.UserService) *UserRouter {
	return &UserRouter{
		UserUsecase: userUsecase,
	}
}

func (ur *UserRouter) Register(w http.ResponseWriter, r *http.Request) {
	var req models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	if err := validator.ValidateUserRegister(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка валидации")
		return
	}

	user, err := ur.UserUsecase.Register(&req)
	if err != nil {
		if err == models.ErrUserAlreadyExists {
			writeErrorResponse(w, http.StatusConflict, err, "Пользователь уже существует")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать пользователя")
		return
	}

	respUser := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"avatar_url": user.AvatarURL,
	}
	writeSuccessResponse(w, http.StatusCreated, respUser)
}

func (ur *UserRouter) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}
	if err := validator.ValidateUserLogin(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка валидации")
		return
	}
	user, token, err := ur.UserUsecase.Login(&req)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			writeErrorResponse(w, http.StatusUnauthorized, err, "Неверные учетные данные")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Ошибка авторизации")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	respUser := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"avatar_url": user.AvatarURL,
	}
	writeSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"user":  respUser,
	})
}

func (ur *UserRouter) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, err, "Токен не найден")
		return
	}
	err = ur.UserUsecase.Logout(cookie.Value)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Ошибка при выходе")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf("Пользователь успешно вышел из системы")
	writeSuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Вы успешно вышли из системы",
	})
}
