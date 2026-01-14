package character

import (
	"encoding/json"
	"errors"
	"net/http"
	"role-helper/internal/delivery/middleware"
	"role-helper/internal/models"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (cr *CharacterRouter) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	var character models.Character

	if err := json.NewDecoder(r.Body).Decode(&character); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}
	character.UserID = user.ID
	createdCharacter, err := cr.CharacterUsecase.Create(&character)
	if err != nil {
		if isValidationError(err) {
			writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка валидации")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать персонажа")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, createdCharacter)
}

func (cr *CharacterRouter) GetCharacters(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}
	characters, err := cr.CharacterUsecase.GetAll(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить список персонажей")
		return
	}

	writeSuccessResponse(w, http.StatusOK, characters)
}

func (cr *CharacterRouter) GetCharacter(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	character, err := cr.CharacterUsecase.FindByID(id, user.ID)
	if err != nil {
		if err == models.ErrCharacterNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Персонаж не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить персонажа")
		return
	}

	writeSuccessResponse(w, http.StatusOK, character)
}

func (cr *CharacterRouter) UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var updateCharacter models.Character
	if err := json.NewDecoder(r.Body).Decode(&updateCharacter); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	updatedCharacter, err := cr.CharacterUsecase.Update(id, user.ID, &updateCharacter)
	if err != nil {
		if err == models.ErrCharacterNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Персонаж не найден")
			return
		}
		if isValidationError(err) {
			writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка валидации")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить персонажа")
		return
	}

	writeSuccessResponse(w, http.StatusOK, updatedCharacter)
}

func (cr *CharacterRouter) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := cr.CharacterUsecase.Delete(id, user.ID)
	if err != nil {
		if err == models.ErrCharacterNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Персонаж не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось удалить персонажа")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cr *CharacterRouter) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, nil, "Требуется авторизация")
		return
	}

	vars := mux.Vars(r)
	characterID, _ := strconv.Atoi(vars["id"])

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		writeErrorResponse(w, http.StatusBadRequest, nil, "Ожидается multipart/form-data")
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка разбора формы")
		return
	}

	file, handler, err := r.FormFile("photo")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Файл 'photo' не найден")
		return
	}
	defer file.Close()

	photoURL, err := cr.CharacterUsecase.UploadPhoto(characterID, user.ID, file, handler.Filename)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Ошибка загрузки фото"

		if strings.Contains(err.Error(), "недопустимый формат") {
			status = http.StatusBadRequest
			msg = "Поддерживаются только JPG, PNG"
		} else if strings.Contains(err.Error(), "превышает допустимый размер") {
			status = http.StatusBadRequest
			msg = "Размер файла не должен превышать 10 МБ"
		} else if strings.Contains(err.Error(), "не принадлежит") {
			status = http.StatusForbidden
			msg = "Персонаж не принадлежит пользователю"
		} else if err == models.ErrCharacterNotFound {
			status = http.StatusNotFound
			msg = "Персонаж не найден"
		}

		writeErrorResponse(w, status, err, msg)
		return
	}

	writeSuccessResponse(w, http.StatusOK, map[string]string{"photo_url": photoURL})
}
