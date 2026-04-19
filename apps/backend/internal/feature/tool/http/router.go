package http

import (
	"encoding/json"
	"net/http"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/internal/infra"
	"tericcabrel/instech/internal/repository"

	"github.com/go-chi/chi/v5"
)

type ToolRouter struct {
	ToolRepository         repository.ToolRepositoryInterface
	RelationshipRepository repository.RelationshipRepositoryInterface
}

func (deps *ToolRouter) Initialize() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var tool usecase.AddToolInput
		err := json.NewDecoder(r.Body).Decode(&tool)

		if err != nil {
			infra.BadRequestError(w, tool)
			return
		}
		addTool := usecase.AddToolUseCase{
			ToolRepository: deps.ToolRepository,
		}

		createdTool, err := addTool.Execute(tool)

		if err != nil {
			if _, ok := err.(domain.ErrInvalidToolCategory); ok {
				infra.BadRequestError(w, err)
				return
			}
			if _, ok := err.(domain.ErrInvalidToolSubType); ok {
				infra.BadRequestError(w, err)
				return
			}
			if _, ok := err.(domain.ErrInvalidToolDevstatus); ok {
				infra.BadRequestError(w, err)
				return
			}
			if _, ok := err.(domain.ErrInvalidField); ok {
				infra.BadRequestError(w, err)
				return
			}
			infra.InternalServerError(w, err, "AddToolUsecase")
			return
		}

		infra.Created(w, createdTool)
	})

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")
		getTool := usecase.GetToolBySlugUseCase{
			ToolRepository: deps.ToolRepository,
		}
		tool, err := getTool.Execute(slug)

		if err != nil {
			if _, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, slug)
				return
			}

			infra.InternalServerError(w, err, "GetToolBySlugUsecase")
			return
		}
		infra.OK(w, tool)
	})

	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")
		deleteTool := usecase.DeleteToolUseCase{
			ToolRepository: deps.ToolRepository,
		}
		err := deleteTool.Execute(slug)
		if err != nil {
			infra.InternalServerError(w, err, "DeleteToolUsecase")
			return
		}
		infra.NoContent(w)
	})

	router.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")

		var tool usecase.UpdateToolInput
		err := json.NewDecoder(r.Body).Decode(&tool)
		if err != nil {
			infra.BadRequestError(w, tool)
			return
		}

		updateTool := usecase.UpdateToolUseCase{
			ToolRepository: deps.ToolRepository,
		}
		updatedTool, err := updateTool.Execute(slug, tool)
		if err != nil {
			if _, ok := err.(domain.ErrInvalidToolCategory); ok {
				infra.BadRequestError(w, err)
				return
			}
			if _, ok := err.(domain.ErrInvalidToolSubType); ok {
				infra.BadRequestError(w, err)
				return
			}
			if _, ok := err.(domain.ErrInvalidToolDevstatus); ok {
				infra.BadRequestError(w, err)
				return
			}
			if _, ok := err.(domain.ErrInvalidField); ok {
				infra.BadRequestError(w, err)
				return
			}
			if _, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, slug)
				return
			}
			infra.InternalServerError(w, err, "UpdateToolUsecase")
			return
		}
		infra.OK(w, updatedTool)
	})

	router.Get("/{id}/alternatives", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")

		getToolAlternatives := usecase.GetToolAlternativesUseCase{
			ToolRepository:         deps.ToolRepository,
			RelationshipRepository: deps.RelationshipRepository,
		}
		alternatives, err := getToolAlternatives.Execute(slug)

		if err != nil {
			if _, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, slug)
				return
			}
			infra.InternalServerError(w, err, "GetToolAlternativesUsecase")
			return
		}

		infra.OK(w, alternatives)
	})

	router.Get("/{id}/graph", func(w http.ResponseWriter, r *http.Request) {
		var result map[string]any = map[string]any{
			"message": "Tool graph for " + chi.URLParam(r, "id"),
		}
		infra.OK(w, result)
	})

	return router
}
