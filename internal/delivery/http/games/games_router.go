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
	gamesRouter.HandleFunc("/{game_id}/sessions", gr.CreateSession).Methods("POST", "OPTIONS")
	gamesRouter.HandleFunc("/sessions/enter", gr.EnterSession).Methods("POST", "OPTIONS")
	gamesRouter.HandleFunc("/sessions/{session_id}/finish", gr.FinishSession).Methods("POST", "OPTIONS")
	gamesRouter.HandleFunc("/{game_id}/players", gr.GetGamePlayers).Methods("GET")
	gamesRouter.HandleFunc("/sessions/{session_id}/players", gr.GetSessionPlayers).Methods("GET")
	gamesRouter.HandleFunc("/{game_id}/sessions", gr.GetPreviousSessions).Methods("GET")
}
