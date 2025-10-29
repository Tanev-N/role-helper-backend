package character

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"role-helper/internal/models"
)

func (cr *CharacterRouter) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	var character models.Character

	if err := json.NewDecoder(r.Body).Decode(&character); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

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
	characters, err := cr.CharacterUsecase.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить список персонажей")
		return
	}

	writeSuccessResponse(w, http.StatusOK, characters)
}

func (cr *CharacterRouter) GetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	character, err := cr.CharacterUsecase.FindByID(id)
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
	vars := mux.Vars(r)
	id := vars["id"]

	var updateCharacter models.Character
	if err := json.NewDecoder(r.Body).Decode(&updateCharacter); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	updatedCharacter, err := cr.CharacterUsecase.Update(id, &updateCharacter)
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
	vars := mux.Vars(r)
	id := vars["id"]

	err := cr.CharacterUsecase.Delete(id)
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
