package httpserver

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"role-helper/cfg"
	"role-helper/internal/delivery/http/character"
	"role-helper/internal/delivery/middleware"
	"role-helper/internal/repository"
	"role-helper/internal/usecase"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start(config *cfg.Config, db *sql.DB, client *redis.Client) error {
	router := s.setupRoutes(db)

	s.server = &http.Server{
		Addr:    config.HTTPServer.IP + ":" + config.HTTPServer.Port,
		Handler: router,
	}
	s.setupRoutes(db)
	log.Println("Server is running on", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) setupRoutes(db *sql.DB) *mux.Router {
	cr := repository.NewCharacterDB(db)
	cu := usecase.NewCharacterUsecase(cr)

	router := mux.NewRouter()
	router = router.PathPrefix("/").Subrouter()

	router.Use(middleware.CORS)

	characterRout := character.NewCharacterRouter(cu)
	characterRout.SetupCharacterRoutes(router)
	return router
}
