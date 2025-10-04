package character

import (
	"encoding/json"
	"log"
	"net/http"
	"role-helper/internal/models"
	"role-helper/internal/validator"

	"github.com/gorilla/mux"
)

func (cr *CharacterRouter) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	log.Printf("Создание нового персонажа")

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

	log.Printf("Персонаж успешно создан: %s (ID: %s)", createdCharacter.Name, createdCharacter.ID)
	writeSuccessResponse(w, http.StatusCreated, createdCharacter)
}

func (cr *CharacterRouter) GetCharacters(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получение всех персонажей")

	characters, err := cr.CharacterUsecase.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить список персонажей")
		return
	}

	log.Printf("Получено %d персонажей", len(characters))
	writeSuccessResponse(w, http.StatusOK, characters)
}

func (cr *CharacterRouter) GetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Получение персонажа по ID: %s", id)

	character, err := cr.CharacterUsecase.FindByID(id)
	if err != nil {
		if err == validator.ErrCharacterNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Персонаж не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить персонажа")
		return
	}

	log.Printf("Персонаж получен: %s", character.Name)
	writeSuccessResponse(w, http.StatusOK, character)
}

func (cr *CharacterRouter) UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Обновление персонажа ID: %s", id)

	var updateCharacter models.Character
	if err := json.NewDecoder(r.Body).Decode(&updateCharacter); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	updatedCharacter, err := cr.CharacterUsecase.Update(id, &updateCharacter)
	if err != nil {
		if err == validator.ErrCharacterNotFound {
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

	log.Printf("Персонаж успешно обновлен: %s", updatedCharacter.Name)
	writeSuccessResponse(w, http.StatusOK, updatedCharacter)
}

func (cr *CharacterRouter) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Удаление персонажа ID: %s", id)

	err := cr.CharacterUsecase.Delete(id)
	if err != nil {
		if err == validator.ErrCharacterNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Персонаж не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось удалить персонажа")
		return
	}

	log.Printf("Персонаж успешно удален: %s", id)
	w.WriteHeader(http.StatusNoContent)
}
