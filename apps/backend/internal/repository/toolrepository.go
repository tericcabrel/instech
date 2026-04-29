package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"tericcabrel/instech/db/queries"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
)

type ToolRepositoryInterface interface {
	CreateTool(ctx context.Context, tool domain.Tool) (domain.Tool, error)
	DeleteTool(ctx context.Context, slug string) error
	GetToolByID(ctx context.Context, id int) (domain.Tool, error)
	GetToolByIDs(ctx context.Context, ids []int) ([]domain.Tool, error)
	GetToolBySlug(ctx context.Context, slug string) (domain.Tool, error)
	UpdateTool(ctx context.Context, tool domain.Tool) (domain.Tool, error)
}

type ToolRepository struct {
	db *sql.DB
}

func NewToolRepository(db *sql.DB) *ToolRepository {
	return &ToolRepository{db: db}
}

func (t *ToolRepository) GetToolByID(ctx context.Context, id int) (domain.Tool, error) {
	record, err := queries.New(t.db).GetToolByID(ctx, id)
	if err != nil {
		return domain.Tool{}, err
	}
	tool, err := MapToolRecordToTool(record)
	if err != nil {
		return domain.Tool{}, err
	}

	return tool, nil
}

func (t *ToolRepository) GetToolBySlug(ctx context.Context, slug string) (domain.Tool, error) {
	record, err := queries.New(t.db).GetToolBySlug(ctx, slug)
	if err != nil {
		return domain.Tool{}, err
	}

	tool, err := MapToolRecordToTool(record)

	if err != nil {
		return domain.Tool{}, err
	}

	return tool, nil
}

func (t *ToolRepository) CreateTool(ctx context.Context, tool domain.Tool) (domain.Tool, error) {
	params := queries.CreateToolParams{
		Name:        tool.Name,
		Slug:        tool.Slug,
		Category:    tool.Category,
		SubType:     sql.NullString{String: tool.SubType, Valid: true},
		Prolang:     sql.NullString{String: tool.Prolang, Valid: true},
		ReleaseYear: tool.ReleaseYear,
		DevStatus:   sql.NullString{String: tool.DevStatus, Valid: true},
		Details:     sql.NullString{String: tool.Details, Valid: true},
		Website:     sql.NullString{String: tool.Website, Valid: true},
		Github:      sql.NullString{String: tool.Github, Valid: true},
	}

	jsonUseCases, err := json.Marshal(tool.UseCases)
	if err != nil {
		return domain.Tool{}, err
	}
	jsonTags, err := json.Marshal(tool.Tags)
	if err != nil {
		return domain.Tool{}, err
	}

	params.UseCases = string(jsonUseCases)
	params.Tags = string(jsonTags)

	record, err := queries.New(t.db).CreateTool(ctx, params)
	if err != nil {
		return domain.Tool{}, err
	}

	tool, mapErr := MapToolRecordToTool(record)

	if mapErr != nil {
		return domain.Tool{}, mapErr
	}

	return tool, nil
}

func (t *ToolRepository) UpdateTool(ctx context.Context, tool domain.Tool) (domain.Tool, error) {
	params := queries.UpdateToolParams{
		Name:        tool.Name,
		Slug:        tool.Slug,
		Category:    tool.Category,
		SubType:     sql.NullString{String: tool.SubType, Valid: true},
		Prolang:     sql.NullString{String: tool.Prolang, Valid: true},
		ReleaseYear: tool.ReleaseYear,
		DevStatus:   sql.NullString{String: tool.DevStatus, Valid: true},
		Details:     sql.NullString{String: tool.Details, Valid: true},
		Website:     sql.NullString{String: tool.Website, Valid: true},
		Github:      sql.NullString{String: tool.Github, Valid: true},
		Id:          tool.Id,
	}

	jsonUseCases, err := json.Marshal(tool.UseCases)
	if err != nil {
		return domain.Tool{}, err
	}
	jsonTags, err := json.Marshal(tool.Tags)
	if err != nil {
		return domain.Tool{}, err
	}

	params.UseCases = string(jsonUseCases)
	params.Tags = string(jsonTags)

	record, err := queries.New(t.db).UpdateTool(ctx, params)
	if err != nil {
		return domain.Tool{}, err
	}

	tool, mapErr := MapToolRecordToTool(record)
	if mapErr != nil {
		return domain.Tool{}, mapErr
	}

	return tool, nil
}

func (t *ToolRepository) DeleteTool(ctx context.Context, slug string) error {
	return queries.New(t.db).DeleteTool(ctx, slug)
}

func (t *ToolRepository) GetToolByIDs(ctx context.Context, ids []int) ([]domain.Tool, error) {
	records, err := queries.New(t.db).GetToolsByIDs(ctx, ids)

	if err != nil {
		return []domain.Tool{}, err
	}

	var result = make([]domain.Tool, 0)

	for _, r := range records {
		tool, err := MapToolRecordToTool(r)
		if err != nil {
			return []domain.Tool{}, common.ErrResourceNotFound{Id: strconv.Itoa(r.Id), Message: "The tool was not found"}
		} else {
			result = append(result, tool)
		}
	}

	return result, nil
}
