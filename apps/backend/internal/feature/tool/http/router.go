package http

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

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

	router.Get("/query", func(w http.ResponseWriter, r *http.Request) {
		keyword := strings.TrimSpace(r.URL.Query().Get("q"))

		searchTools := usecase.SearchToolsUseCase{
			ToolRepository: deps.ToolRepository,
		}
		results, err := searchTools.Execute(keyword)
		if err != nil {
			httprouter.InternalServerError(w, err, "SearchToolsUsecase")

			return
		}

		httprouter.OK(w, results)
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
		slug := chi.URLParam(r, "id")
		depth := 1
		depthParam := strings.TrimSpace(r.URL.Query().Get("depth"))
		if depthParam != "" {
			parsedDepth, err := strconv.Atoi(depthParam)
			if err == nil && parsedDepth > 0 {
				depth = int(math.Min(float64(parsedDepth), float64(usecase.MAX_GRAPH_DEPTH)))
			}
		}

		layoutMode := strings.TrimSpace(r.URL.Query().Get("layoutMode"))
		if layoutMode == "" || !usecase.IsLayoutModeValid(layoutMode) {
			layoutMode = usecase.GRAPH_LAYOUT_MODE_CHRONOLOGICAL
		}

		kinds := r.URL.Query()["kinds"]
		validKinds := make([]string, 0, len(kinds))
		for _, kind := range kinds {
			if domain.IsKindValid(kind) {
				validKinds = append(validKinds, kind)
			}
		}

		getToolGraph := usecase.GetToolGraphUseCase{
			ToolRepository:         deps.ToolRepository,
			RelationshipRepository: deps.RelationshipRepository,
		}
		graph, err := getToolGraph.Execute(slug, usecase.GetToolGraphInput{
			Depth:      depth,
			Kinds:      validKinds,
			LayoutMode: layoutMode,
		})

		if err != nil {
			if _, ok := err.(common.ErrResourceNotFound); ok {
				httprouter.NotFoundError(w, slug)

				return
			}
			httprouter.InternalServerError(w, err, "GetToolGraphUsecase")

			return
		}

		httprouter.OK(w, graph)
	})

	return router
}
