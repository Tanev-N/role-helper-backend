package httpserver

import (
	"database/sql"
	"log"
	"net/http"
	"role-helper/cfg"
	"role-helper/internal/delivery/http/armor"
	"role-helper/internal/delivery/http/character"
	games "role-helper/internal/delivery/http/games"
	"role-helper/internal/delivery/http/user"
	"role-helper/internal/delivery/http/weapon"
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

	gameRepoInterface := repository.NewGameRepository(db)
	gameRepo, ok := gameRepoInterface.(*repository.GameRepository)
	if !ok {
		log.Fatalf("unexpected type for game repository: %T", gameRepoInterface)
	}
	gu := usecase.NewGameUseCase(gameRepo, ur, cr)

	ar := repository.NewArmorRepository(db)
	au := usecase.NewArmorUsecase(ar)

	wr := repository.NewWeaponRepository(db)
	wu := usecase.NewWeaponUsecase(wr)

	router := mux.NewRouter()

	openapiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "role-helper-api.yaml")
	})
	router.Path("/openapi.yaml").Handler(middleware.CORS(openapiHandler)).Methods("GET", "OPTIONS")

	docsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>DnD Characters API</title>
            <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui.css">
            <style>
                body { margin: 0; }
                #swagger-ui { padding: 20px; }
                .swagger-ui .topbar { display: none; }
            </style>
        </head>
        <body>
            <div id="swagger-ui"></div>
            <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
            <script>
                window.onload = function() {
                    const protocol = window.location.protocol;
                    const host = window.location.host;
                    const openapiUrl = protocol + '//' + host + '/openapi.yaml';
                    window.ui = SwaggerUIBundle({
                        url: openapiUrl,
                        dom_id: '#swagger-ui',
                        presets: [
                            SwaggerUIBundle.presets.apis,
                            SwaggerUIBundle.presets.standalone
                        ],
                        requestInterceptor: (request) => {
                            return request;
                        }
                    });
                }
            </script>
        </body>
        </html>`))
	})
	router.Path("/docs").Handler(middleware.CORS(docsHandler)).Methods("GET", "OPTIONS")

	api := router.PathPrefix("/api").Subrouter()

	api.Use(middleware.CORS)
	api.Use(middleware.Auth(uu))

	characterRout := character.NewCharacterRouter(cu)
	characterRout.SetupCharacterRoutes(api)

	gameRout := games.NewGameRouter(gu)
	gameRout.SetupRoutes(api)

	userRout := user.NewUserRouter(uu)
	userRout.SetupRoutes(api)

	armorRout := armor.NewArmorRouter(au)
	armorRout.SetupArmorRoutes(api)

	weaponRout := weapon.NewWeaponRouter(wu)
	weaponRout.SetupWeaponRoutes(api)

	return router
}
