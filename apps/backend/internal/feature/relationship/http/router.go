package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/relationship/usecase"
	"tericcabrel/instech/internal/infra"
	"tericcabrel/instech/internal/repository"

	"github.com/go-chi/chi/v5"
)

type RelationshipRouter struct {
	RelationshipRepository repository.RelationshipRepositoryInterface
	ToolRepository         repository.ToolRepositoryInterface
}

func (deps *RelationshipRouter) Initialize() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var input usecase.CreateRelationshipInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			infra.BadRequestError(w, input)
			return
		}
		createRelationship := usecase.CreateRelationshipUseCase{
			RelationshipRepository: deps.RelationshipRepository,
			ToolRepository:         deps.ToolRepository,
		}
		relationship, err := createRelationship.Execute(input)
		if err != nil {
			if errToolNotFound, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, errToolNotFound.Id)
				return
			}
			if errInvalidRelationshipKind, ok := err.(domain.ErrInvalidRelationshipKind); ok {
				infra.BadRequestError(w, map[string]string{
					"message": errInvalidRelationshipKind.Message,
					"kind":    errInvalidRelationshipKind.Kind,
				})
				return
			}
			if errInvalidField, ok := err.(domain.ErrInvalidField); ok {
				infra.BadRequestError(w, errInvalidField)
				return
			}
			infra.InternalServerError(w, err, "CreateRelationshipUseCase")
			return
		}
		infra.Created(w, relationship)
	})

	router.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, castErr := strconv.Atoi(chi.URLParam(r, "id"))
		if castErr != nil {
			infra.BadRequestError(w, map[string]string{
				"message": "Invalid relationship ID",
			})
			return
		}

		var input usecase.UpdateRelationshipInput
		parseErr := json.NewDecoder(r.Body).Decode(&input)
		if parseErr != nil {
			infra.BadRequestError(w, input)
			return
		}

		updateRelationship := usecase.UpdateRelationshipUseCase{
			RelationshipRepository: deps.RelationshipRepository,
			ToolRepository:         deps.ToolRepository,
		}
		updatedRelationship, err := updateRelationship.Execute(id, input)

		if err != nil {
			if errToolNotFound, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, errToolNotFound.Id)
				return
			}
			if errInvalidRelationshipKind, ok := err.(domain.ErrInvalidRelationshipKind); ok {
				infra.BadRequestError(w, map[string]string{
					"message": errInvalidRelationshipKind.Message,
					"kind":    errInvalidRelationshipKind.Kind,
				})
				return
			}
			if errInvalidField, ok := err.(domain.ErrInvalidField); ok {
				infra.BadRequestError(w, errInvalidField)
				return
			}
			infra.InternalServerError(w, err, "UpdateRelationshipUseCase")
			return
		}
		infra.OK(w, updatedRelationship)
	})

	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, castErr := strconv.Atoi(chi.URLParam(r, "id"))
		if castErr != nil {
			infra.BadRequestError(w, map[string]string{
				"message": "Invalid relationship Id",
			})
			return
		}

		deleteRelationship := usecase.DeleteRelationshipUseCase{
			RelationshipRepository: deps.RelationshipRepository,
		}
		err := deleteRelationship.Execute(id)
		if err != nil {
			if errResourceNotFound, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, errResourceNotFound.Id)
			} else {
				infra.InternalServerError(w, err, "DeleteRelationshipUseCase")
			}
			return
		}
		infra.NoContent(w)
	})

	router.Get("/query", func(w http.ResponseWriter, r *http.Request) {
		toolIdParam := r.URL.Query().Get("tool_id")
		kindParam := r.URL.Query().Get("kind")
		cursorParam := r.URL.Query().Get("cursor")
		limitParam := r.URL.Query().Get("limit")

		var toolId int
		var kind string
		var cursor int64
		var castErr error

		if toolIdParam != "" {
			toolId, castErr = strconv.Atoi(toolIdParam)
			if castErr != nil {
				infra.BadRequestError(w, map[string]string{
					"message": "Invalid tool Id",
				})
				return
			}
		}

		if kindParam != "" {
			kind = kindParam
			if !domain.IsKindValid(kind) {
				infra.BadRequestError(w, map[string]string{
					"message": "Invalid kind",
					"kind":    kind,
				})
				return
			}
		}

		if cursorParam != "" {
			cursor, castErr = strconv.ParseInt(cursorParam, 10, 64)
			if castErr != nil {
				infra.BadRequestError(w, map[string]string{
					"message": "Invalid cursor",
				})
				return
			}
		}

		var limit int
		if limitParam != "" {
			limit, castErr = strconv.Atoi(limitParam)
			if castErr != nil {
				infra.BadRequestError(w, map[string]string{
					"message": "Invalid limit",
				})
				return
			}
		}

		getRelationships := usecase.GetRelationshipsUseCase{
			RelationshipRepository: deps.RelationshipRepository,
			ToolRepository:         deps.ToolRepository,
		}
		results, err := getRelationships.Execute(usecase.GetRelationshipsUseCaseParams{
			Cursor: cursor,
			ToolId: toolId,
			Kind:   kind,
			Limit:  limit,
		})

		if err != nil {
			if errResourceNotFound, ok := err.(common.ErrResourceNotFound); ok {
				infra.NotFoundError(w, errResourceNotFound.Id)
				return
			}
			infra.InternalServerError(w, err, "GetRelationshipsUsecase")
			return
		}

		infra.OK(w, results)
	})

	return router
}
