package user

import (
	"github.com/gorilla/mux"
)

func (ur *UserRouter) SetupRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/api/auth").Subrouter()

	userRouter.HandleFunc("/register", ur.Register).Methods("POST", "OPTIONS")
	userRouter.HandleFunc("/login", ur.Login).Methods("POST", "OPTIONS")
	userRouter.HandleFunc("/logout", ur.Logout).Methods("POST", "OPTIONS")
}
