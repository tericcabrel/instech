package http

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/relationship/usecase"
	"tericcabrel/instech/internal/infra"
	"tericcabrel/instech/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func InitializeRelationshipRouter(relationshipRepository repository.RelationshipRepositoryInterface, toolRepository repository.ToolRepositoryInterface) *chi.Mux {
	relationshipRouter := chi.NewRouter()

	relationshipRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var input usecase.CreateRelationshipInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			infra.BadRequestError(w, input)
			return
		}
		relationship, err := usecase.CreateRelationshipUseCase(relationshipRepository, toolRepository, input)
		if err != nil {
			if errToolNotFound, ok := err.(common.ErrResourceNotFound); ok {
				infra.BadRequestError(w, map[string]string{
					"message": errToolNotFound.Message,
					"slug":    errToolNotFound.Id,
				})
				return
			}
			if errInvalidRelationshipKind, ok := err.(common.ErrInvalidRelationshipKind); ok {
				infra.BadRequestError(w, map[string]string{
					"message": errInvalidRelationshipKind.Message,
					"kind":    errInvalidRelationshipKind.Kind,
				})
				return
			}
			infra.InternalServerError(w, err, "CreateRelationshipUseCase")
			return
		}
		infra.Created(w, relationship)
	})

	relationshipRouter.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		updatedRelationship, err := usecase.UpdateRelationshipUseCase(relationshipRepository, toolRepository, id, input)

		if err != nil {
			if errToolNotFound, ok := err.(common.ErrResourceNotFound); ok {
				infra.BadRequestError(w, map[string]string{
					"message": errToolNotFound.Message,
					"slug":    errToolNotFound.Id,
				})
				return
			}
			if errInvalidRelationshipKind, ok := err.(common.ErrInvalidRelationshipKind); ok {
				infra.BadRequestError(w, map[string]string{
					"message": errInvalidRelationshipKind.Message,
					"kind":    errInvalidRelationshipKind.Kind,
				})
				return
			}
			infra.InternalServerError(w, err, "UpdateRelationshipUseCase")
			return
		}
		infra.OK(w, updatedRelationship)
	})

	relationshipRouter.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, castErr := strconv.Atoi(chi.URLParam(r, "id"))
		if castErr != nil {
			infra.BadRequestError(w, map[string]string{
				"message": "Invalid relationship Id",
			})
			return
		}

		err := usecase.DeleteRelationshipUseCase(relationshipRepository, id)
		if err != nil {
			if errResourceNotFound, ok := err.(common.ErrResourceNotFound); ok {
				log.Log().Err(errResourceNotFound).Msg("Relationship not found")
			} else {
				infra.InternalServerError(w, err, "DeleteRelationshipUseCase")
				return
			}
		}
		infra.NoContent(w)
	})

	relationshipRouter.Get("/query", func(w http.ResponseWriter, r *http.Request) {
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
			if !slices.Contains(domain.RELATIONSHIP_KINDS, kind) {
				infra.BadRequestError(w, map[string]string{
					"message": "Invalid kind",
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

		results, err := usecase.GetRelationshipsUsecase(relationshipRepository, toolRepository, usecase.GetRelationshipsUseCaseParams{
			Cursor: cursor,
			ToolId: toolId,
			Kind:   kind,
			Limit:  limit,
		})

		if err != nil {
			infra.InternalServerError(w, err, "GetRelationshipsUsecase")
			return
		}

		infra.OK(w, results)
	})

	return relationshipRouter
}
