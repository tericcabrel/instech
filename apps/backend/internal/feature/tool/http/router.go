package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/internal/infra"
	"tericcabrel/instech/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func InitializeToolRouter(toolRepository repository.ToolRepositoryInterface) *chi.Mux {
	toolRouter := chi.NewRouter()

	toolRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var tool usecase.AddToolInput
		err := json.NewDecoder(r.Body).Decode(&tool)

		if err != nil {
			infra.BadRequestError(w, tool)
			return
		}

		createdTool, err := usecase.AddToolUsecase(toolRepository, tool)

		log.Info().Msgf("Created tool: %+v", createdTool)

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
			fmt.Printf("Error: %+v\n %t", err, errors.Is(err, usecase.ErrToolNotFound{}))
			if _, ok := err.(usecase.ErrToolNotFound); ok {
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
			if _, ok := err.(usecase.ErrToolNotFound); ok {
				infra.NotFoundError(w, slug)
				return
			}
			infra.InternalServerError(w, err, "UpdateToolUsecase")
			return
		}
		infra.OK(w, updatedTool)
	})

	toolRouter.Get("/{id}/alternatives", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Tool alternatives for %s", chi.URLParam(r, "id"))))
	})
	toolRouter.Get("/{id}/similar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Tool similar to %s", chi.URLParam(r, "id"))))
	})
	toolRouter.Get("/{id}/relationships", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Tool relationships for %s", chi.URLParam(r, "id"))))
	})

	toolRouter.Get("/{id}/graph", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Tool graph for %s", chi.URLParam(r, "id"))))
	})

	return toolRouter
}
