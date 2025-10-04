package character

import (
	"role-helper/internal/models"

	"github.com/gorilla/mux"
)

type CharacterRouter struct {
	CharacterUsecase models.CharacterService
}

func NewCharacterRouter(cs models.CharacterService) *CharacterRouter {
	return &CharacterRouter{CharacterUsecase: cs}
}

func (cr *CharacterRouter) SetupCharacterRoutes(mux *mux.Router) {
	charactersRouter := mux.PathPrefix("/characters").Subrouter()

	charactersRouter.HandleFunc("", cr.CreateCharacter).Methods("POST", "OPTIONS")
	charactersRouter.HandleFunc("", cr.GetCharacters).Methods("GET")
	charactersRouter.HandleFunc("/{id}", cr.GetCharacter).Methods("GET")
	charactersRouter.HandleFunc("/{id}", cr.UpdateCharacter).Methods("PUT", "OPTIONS")
	charactersRouter.HandleFunc("/{id}", cr.DeleteCharacter).Methods("DELETE", "OPTIONS")
}
