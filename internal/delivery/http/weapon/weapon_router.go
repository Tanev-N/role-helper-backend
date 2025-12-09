package weapon

import (
	"role-helper/internal/models"

	"github.com/gorilla/mux"
)

type WeaponRouter struct {
	WeaponUsecase models.WeaponService
}

func NewWeaponRouter(ws models.WeaponService) *WeaponRouter {
	return &WeaponRouter{WeaponUsecase: ws}
}

func (wr *WeaponRouter) SetupWeaponRoutes(mux *mux.Router) {
	weaponRouter := mux.PathPrefix("/weapon").Subrouter()

	weaponRouter.HandleFunc("", wr.CreateWeapon).Methods("POST", "OPTIONS")
	weaponRouter.HandleFunc("", wr.GetWeapons).Methods("GET")
	weaponRouter.HandleFunc("/{id}", wr.GetWeapon).Methods("GET")
	weaponRouter.HandleFunc("/{id}", wr.UpdateWeapon).Methods("PUT", "OPTIONS")
	weaponRouter.HandleFunc("/{id}", wr.DeleteWeapon).Methods("DELETE", "OPTIONS")
}
