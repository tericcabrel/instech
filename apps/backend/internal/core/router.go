package core

import (
	"net/http"
	relationshiphttp "tericcabrel/instech/internal/feature/relationship/http"
	toolhttp "tericcabrel/instech/internal/feature/tool/http"
	"tericcabrel/instech/internal/infra/httprouter"
	"tericcabrel/instech/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type HTTPRouter struct {
	ToolRepository         repository.ToolRepositoryInterface
	RelationshipRepository repository.RelationshipRepositoryInterface
}

func (router *HTTPRouter) Initialize() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://instech.com", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		httprouter.OK(w, map[string]string{"message": "Hello from Instech"})
	})

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		httprouter.OK(w, map[string]string{"message": "Search results for " + r.URL.Query().Get("q")})
	})

	toolRouter := toolhttp.ToolRouter{
		ToolRepository:         router.ToolRepository,
		RelationshipRepository: router.RelationshipRepository,
	}

	r.Mount("/tools", toolRouter.Initialize())

	relationshipRouter := relationshiphttp.RelationshipRouter{
		RelationshipRepository: router.RelationshipRepository,
		ToolRepository:         router.ToolRepository,
	}

	r.Mount("/relationships", relationshipRouter.Initialize())

	return r
}
