package games

import (
	"role-helper/internal/models"

	"github.com/gorilla/mux"
)

type GameRouter struct {
	GameUsecase models.GameService
}

func NewGameRouter(gs models.GameService) *GameRouter {
	return &GameRouter{GameUsecase: gs}
}

func (gr *GameRouter) SetupRoutes(router *mux.Router) {
	gamesRouter := router.PathPrefix("/games").Subrouter()

	gamesRouter.HandleFunc("", gr.CreateGame).Methods("POST", "OPTIONS")
	gamesRouter.HandleFunc("", gr.GetGames).Methods("GET")
	gamesRouter.HandleFunc("/{game_id:[0-9]+}/sessions", gr.CreateSession).Methods("POST", "OPTIONS")
	gamesRouter.HandleFunc("/sessions/enter", gr.EnterSession).Methods("POST", "OPTIONS")
	gamesRouter.HandleFunc("/sessions/{session_id:[0-9]+}/finish", gr.FinishSession).Methods("POST", "OPTIONS")
	gamesRouter.HandleFunc("/{game_id:[0-9]+}/players", gr.GetGamePlayers).Methods("GET")
}
