package armor

import (
	"encoding/json"
	"errors"
	"net/http"
	"role-helper/internal/delivery/middleware"
	"role-helper/internal/models"
	"strconv"

	"github.com/gorilla/mux"
)

func (ar *ArmorRouter) CreateArmor(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	var armor models.Armor
	if err := json.NewDecoder(r.Body).Decode(&armor); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}
	armor.UserID = user.ID

	createdArmor, err := ar.ArmorUsecase.Create(&armor)
	if err != nil {
		if isValidationError(err) {
			writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка валидации")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать броню")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, createdArmor)
}

func (ar *ArmorRouter) GetArmors(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	armors, err := ar.ArmorUsecase.GetAll(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить список брони")
		return
	}

	writeSuccessResponse(w, http.StatusOK, armors)
}

func (ar *ArmorRouter) GetArmor(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	armor, err := ar.ArmorUsecase.GetByID(id, user.ID)
	if err != nil {
		if err == models.ErrArmorNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Броня не найдена")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить броню")
		return
	}

	writeSuccessResponse(w, http.StatusOK, armor)
}

func (ar *ArmorRouter) UpdateArmor(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var updateArmor models.Armor
	if err := json.NewDecoder(r.Body).Decode(&updateArmor); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	updatedArmor, err := ar.ArmorUsecase.Update(id, user.ID, &updateArmor)
	if err != nil {
		if err == models.ErrArmorNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Броня не найдена")
			return
		}
		if isValidationError(err) {
			writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка валидации")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить броню")
		return
	}

	writeSuccessResponse(w, http.StatusOK, updatedArmor)
}

func (ar *ArmorRouter) DeleteArmor(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := ar.ArmorUsecase.Delete(id, user.ID)
	if err != nil {
		if err == models.ErrArmorNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Броня не найдена")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось удалить броню")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
