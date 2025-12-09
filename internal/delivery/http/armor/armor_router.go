package armor

import (
	"role-helper/internal/models"

	"github.com/gorilla/mux"
)

type ArmorRouter struct {
	ArmorUsecase models.ArmorService
}

func NewArmorRouter(as models.ArmorService) *ArmorRouter {
	return &ArmorRouter{ArmorUsecase: as}
}

func (ar *ArmorRouter) SetupArmorRoutes(mux *mux.Router) {
	armorRouter := mux.PathPrefix("/armor").Subrouter()

	armorRouter.HandleFunc("", ar.CreateArmor).Methods("POST", "OPTIONS")
	armorRouter.HandleFunc("", ar.GetArmors).Methods("GET")
	armorRouter.HandleFunc("/{id}", ar.GetArmor).Methods("GET")
	armorRouter.HandleFunc("/{id}", ar.UpdateArmor).Methods("PUT", "OPTIONS")
	armorRouter.HandleFunc("/{id}", ar.DeleteArmor).Methods("DELETE", "OPTIONS")
}
