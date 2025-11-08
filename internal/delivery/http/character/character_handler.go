package character

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"role-helper/internal/delivery/middleware"
	"role-helper/internal/models"
	"strconv"
)

func (cr *CharacterRouter) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
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

	characters, err := cr.CharacterUsecase.GetAll(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить список персонажей")
		return
	}

	writeSuccessResponse(w, http.StatusOK, characters)
}

func (cr *CharacterRouter) GetCharacter(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
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
