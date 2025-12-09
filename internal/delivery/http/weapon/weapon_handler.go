package weapon

import (
	"encoding/json"
	"errors"
	"net/http"
	"role-helper/internal/delivery/middleware"
	"role-helper/internal/models"
	"strconv"

	"github.com/gorilla/mux"
)

func (wr *WeaponRouter) CreateWeapon(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	var weapon models.Weapon
	if err := json.NewDecoder(r.Body).Decode(&weapon); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}
	weapon.UserID = user.ID

	createdWeapon, err := wr.WeaponUsecase.Create(&weapon)
	if err != nil {
		if isValidationError(err) {
			writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка валидации")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать оружие")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, createdWeapon)
}

func (wr *WeaponRouter) GetWeapons(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	weapons, err := wr.WeaponUsecase.GetAll(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить список оружия")
		return
	}

	writeSuccessResponse(w, http.StatusOK, weapons)
}

func (wr *WeaponRouter) GetWeapon(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	weapon, err := wr.WeaponUsecase.GetByID(id, user.ID)
	if err != nil {
		if err == models.ErrWeaponNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Оружие не найдено")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить оружие")
		return
	}

	writeSuccessResponse(w, http.StatusOK, weapon)
}

func (wr *WeaponRouter) UpdateWeapon(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var updateWeapon models.Weapon
	if err := json.NewDecoder(r.Body).Decode(&updateWeapon); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	updatedWeapon, err := wr.WeaponUsecase.Update(id, user.ID, &updateWeapon)
	if err != nil {
		if err == models.ErrWeaponNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Оружие не найдено")
			return
		}
		if isValidationError(err) {
			writeErrorResponse(w, http.StatusBadRequest, err, "Ошибка валидации")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить оружие")
		return
	}

	writeSuccessResponse(w, http.StatusOK, updatedWeapon)
}

func (wr *WeaponRouter) DeleteWeapon(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Необходима авторизация")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := wr.WeaponUsecase.Delete(id, user.ID)
	if err != nil {
		if err == models.ErrWeaponNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Оружие не найдено")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось удалить оружие")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
