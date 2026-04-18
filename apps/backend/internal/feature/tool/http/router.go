package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/internal/infra"
	"tericcabrel/instech/internal/repository"

	"github.com/go-chi/chi/v5"
)

func InitializeToolRouter(toolRepository repository.ToolRepositoryInterface, relationshipRepository repository.RelationshipRepositoryInterface) *chi.Mux {
	toolRouter := chi.NewRouter()

	toolRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var tool usecase.AddToolInput
		err := json.NewDecoder(r.Body).Decode(&tool)

		if err != nil {
			infra.BadRequestError(w, tool)
			return
		}

		createdTool, err := usecase.AddToolUsecase(toolRepository, tool)

		if err != nil {
			infra.InternalServerError(w, err, "AddToolUsecase")
			return
		}

		infra.Created(w, createdTool)
	})

	toolRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")
		tool, err := usecase.GetToolBySlugUsecase(toolRepository, slug)

		if err != nil {
			fmt.Printf("Error: %+v\n %t", err, errors.Is(err, common.ErrResourceNotFound{}))
			if _, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, slug)
				return
			}

			infra.InternalServerError(w, err, "GetToolBySlugUsecase")
			return
		}
		infra.OK(w, tool)
	})

	toolRouter.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")
		err := usecase.DeleteToolUsecase(toolRepository, slug)
		if err != nil {
			infra.InternalServerError(w, err, "DeleteToolUsecase")
			return
		}
		infra.NoContent(w)
	})

	toolRouter.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")

		var tool usecase.UpdateToolInput
		err := json.NewDecoder(r.Body).Decode(&tool)
		if err != nil {
			infra.BadRequestError(w, tool)
			return
		}

		updatedTool, err := usecase.UpdateToolUsecase(toolRepository, slug, tool)
		if err != nil {
			if _, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, slug)
				return
			}
			infra.InternalServerError(w, err, "UpdateToolUsecase")
			return
		}
		infra.OK(w, updatedTool)
	})

	toolRouter.Get("/{id}/alternatives", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")

		alternatives, err := usecase.GetToolAlternativesUsecase(toolRepository, relationshipRepository, slug)

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

	toolRouter.Get("/{id}/similar", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")
		similar, err := usecase.GetSimilarToolUsecase(toolRepository, relationshipRepository, slug)
		if err != nil {
			if _, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, slug)
				return
			}
			infra.InternalServerError(w, err, "GetSimilarToolUsecase")
			return
		}
		infra.OK(w, similar)
	})

	toolRouter.Get("/{id}/graph", func(w http.ResponseWriter, r *http.Request) {
		var result map[string]any = map[string]any{
			"message": "Tool graph for " + chi.URLParam(r, "id"),
		}
		infra.OK(w, result)
	})

	return toolRouter
}
