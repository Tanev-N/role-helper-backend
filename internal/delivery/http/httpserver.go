package httpserver

import (
	"database/sql"
	"log"
	"net/http"
	"role-helper/cfg"
	"role-helper/internal/delivery/http/character"
	"role-helper/internal/delivery/http/user"
	"role-helper/internal/delivery/middleware"
	"role-helper/internal/repository"
	"role-helper/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start(config *cfg.Config, db *sql.DB, client *redis.Client) error {
	router := s.setupRoutes(db, client)

	s.server = &http.Server{
		Addr:    config.HTTPServer.IP + ":" + config.HTTPServer.Port,
		Handler: router,
	}
	log.Println("Server is running on", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) setupRoutes(db *sql.DB, client *redis.Client) *mux.Router {
	cr := repository.NewCharacterRepository(db)
	cu := usecase.NewCharacterUsecase(cr)

	ur := repository.NewUserRepository(db)
	uu := usecase.NewUserUsecase(ur, client)

	router := mux.NewRouter()
	router = router.PathPrefix("/api").Subrouter()

	router.Use(middleware.CORS)
	router.Use(middleware.Auth(uu))

	characterRout := character.NewCharacterRouter(cu)
	characterRout.SetupCharacterRoutes(router)

	userRout := user.NewUserRouter(uu)
	userRout.SetupRoutes(router)

	return router
}
