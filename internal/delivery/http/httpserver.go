package httpserver

import (
	"log"
	"net/http"
	"role-helper/internal/delivery/http/character"
	"role-helper/internal/middleware"
	"role-helper/internal/repository"
	"role-helper/internal/usecase"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start() error {
	router := s.setupRoutes()

	s.server = &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	s.setupRoutes()
	log.Println("Server is running on", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) setupRoutes() *mux.Router {
	cr := repository.NewCharacterMemory()
	cu := usecase.NewCharacterUsecase(cr)

	router := mux.NewRouter()
	router = router.PathPrefix("/").Subrouter()

	router.Use(middleware.CORS)

	characterRout := character.NewCharacterRouter(cu)
	characterRout.SetupCharacterRoutes(router)
	return router
}
