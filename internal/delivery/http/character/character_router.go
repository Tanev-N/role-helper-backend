package character

import (
	"github.com/gorilla/mux"
	"role-helper/internal/models"
)

type CharacterRouter struct {
	CharacterUsecase models.CharacterService
}

func NewCharacterRouter(cs models.CharacterService) *CharacterRouter {
	return &CharacterRouter{CharacterUsecase: cs}
}

func (cr *CharacterRouter) SetupCharacterRoutes(mux *mux.Router) {
	charactersRouter := mux.PathPrefix("/characters").Subrouter()

	charactersRouter.HandleFunc("", cr.CreateCharacter).Methods("POST")
	charactersRouter.HandleFunc("", cr.GetCharacters).Methods("GET")
	charactersRouter.HandleFunc("/{id}", cr.GetCharacter).Methods("GET")
	charactersRouter.HandleFunc("/{id}", cr.UpdateCharacter).Methods("PUT")
	charactersRouter.HandleFunc("/{id}", cr.DeleteCharacter).Methods("DELETE")
}
