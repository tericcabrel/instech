package http

import (
	"encoding/json"
	"net/http"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/internal/infra/httprouter"
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
			httprouter.BadRequestError(w, err.Error())

			return
		}
		addTool := usecase.AddToolUseCase{
			ToolRepository: deps.ToolRepository,
		}

		createdTool, err := addTool.Execute(tool)

		if err != nil {
			switch err.(type) {
			case domain.ErrInvalidToolCategory:
				httprouter.BadRequestError(w, err)

				return
			case domain.ErrInvalidToolSubType:
				httprouter.BadRequestError(w, err)

				return
			case domain.ErrInvalidToolDevStatus:
				httprouter.BadRequestError(w, err)

				return
			case domain.ErrInvalidField:
				httprouter.UnprocessableEntityError(w, err)

				return
			case common.ErrResourceAlreadyExists:
				httprouter.BadRequestError(w, err)

				return
			}

			httprouter.InternalServerError(w, err, "AddToolUsecase")

			return
		}

		httprouter.Created(w, createdTool)
	})

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")
		getTool := usecase.GetToolBySlugUseCase{
			ToolRepository: deps.ToolRepository,
		}
		tool, err := getTool.Execute(slug)

		if err != nil {
			if _, ok := err.(common.ErrResourceNotFound); ok {
				httprouter.NotFoundError(w, slug)

				return
			}

			httprouter.InternalServerError(w, err, "GetToolBySlugUsecase")

			return
		}
		httprouter.OK(w, tool)
	})

	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")
		deleteTool := usecase.DeleteToolUseCase{
			ToolRepository: deps.ToolRepository,
		}
		err := deleteTool.Execute(slug)
		if err != nil {
			httprouter.InternalServerError(w, err, "DeleteToolUsecase")

			return
		}
		httprouter.NoContent(w)
	})

	router.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "id")

		var tool usecase.PatchToolInput
		err := json.NewDecoder(r.Body).Decode(&tool)
		if err != nil {
			httprouter.BadRequestError(w, err.Error())

			return
		}

		updateTool := usecase.PatchToolUseCase{
			ToolRepository: deps.ToolRepository,
		}
		updatedTool, err := updateTool.Execute(slug, tool)
		if err != nil {
			switch err.(type) {
			case domain.ErrInvalidToolCategory:
				httprouter.BadRequestError(w, err)

				return
			case domain.ErrInvalidToolSubType:
				httprouter.BadRequestError(w, err)

				return
			case domain.ErrInvalidToolDevStatus:
				httprouter.BadRequestError(w, err)

				return
			case domain.ErrInvalidField:
				httprouter.UnprocessableEntityError(w, err)

				return
			case common.ErrResourceNotFound:
				httprouter.NotFoundError(w, slug)

				return
			}

			httprouter.InternalServerError(w, err, "PatchToolUsecase")

			return
		}
		httprouter.OK(w, updatedTool)
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
				httprouter.NotFoundError(w, slug)

				return
			}
			httprouter.InternalServerError(w, err, "GetToolAlternativesUsecase")

			return
		}

		httprouter.OK(w, alternatives)
	})

	router.Get("/{id}/graph", func(w http.ResponseWriter, r *http.Request) {
		var result = map[string]any{
			"message": "Tool graph for " + chi.URLParam(r, "id"),
		}
		httprouter.OK(w, result)
	})

	return router
}
