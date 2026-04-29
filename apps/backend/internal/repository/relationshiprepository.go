package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"tericcabrel/instech/db/queries"
	"tericcabrel/instech/internal/domain"
	"time"
)

const DEFAULT_LIMIT = 25
const MAX_LIMIT = 100

type PaginatedRelationshipsResult struct {
	Relationships []domain.Relationship
	TotalCount    int64
	ItemsCount    int64
	NextCursor    int64
}

type GetRelationshipsAllParams struct {
	Kind   string
	Cursor int64
	ToolId int
	Limit  int
}

type RelationshipRepositoryInterface interface {
	CreateRelationship(ctx context.Context, relationship domain.Relationship) (domain.Relationship, error)
	GetRelationshipByID(ctx context.Context, id int) (domain.Relationship, error)
	DeleteRelationship(ctx context.Context, id int) error
	GetRelationshipsByToolID(ctx context.Context, toolID int) ([]domain.Relationship, error)
	UpdateRelationship(ctx context.Context, relationship domain.Relationship) (domain.Relationship, error)
	GetToolAlternatives(ctx context.Context, toolID int) ([]domain.Relationship, error)
	GetRelationshipsAll(ctx context.Context, params GetRelationshipsAllParams) (PaginatedRelationshipsResult, error)
}

type RelationshipRepository struct {
	db *sql.DB
}

func NewRelationshipRepository(db *sql.DB) *RelationshipRepository {
	return &RelationshipRepository{db: db}
}

func (r *RelationshipRepository) CreateRelationship(ctx context.Context, relationship domain.Relationship) (domain.Relationship, error) {
	params := queries.CreateRelationshipParams{
		FromToolId: relationship.FromToolID,
		ToToolId:   relationship.ToToolID,
		Kind:       relationship.Kind,
	}

	jsonMetadata, err := json.Marshal(relationship.Metadata)
	if err != nil {
		return domain.Relationship{}, err
	}
	params.Metadata = string(jsonMetadata)

	record, err := queries.New(r.db).CreateRelationship(ctx, params)
	if err != nil {
		return domain.Relationship{}, err
	}

	return MapRelationshipRecordToRelationship(record), nil
}

func (r *RelationshipRepository) GetRelationshipByID(ctx context.Context, id int) (domain.Relationship, error) {
	record, err := queries.New(r.db).GetRelationshipByID(ctx, id)
	if err != nil {
		return domain.Relationship{}, err
	}

	return MapRelationshipRecordToRelationship(record), nil
}

func (r *RelationshipRepository) DeleteRelationship(ctx context.Context, id int) error {
	return queries.New(r.db).DeleteRelationship(ctx, id)
}

func (r *RelationshipRepository) GetRelationshipsByToolID(ctx context.Context, toolID int) ([]domain.Relationship, error) {
	params := queries.GetRelationshipsByToolIDParams{
		FromToolId: toolID,
		ToToolId:   toolID,
	}

	records, err := queries.New(r.db).GetRelationshipsByToolID(ctx, params)
	if err != nil {
		return []domain.Relationship{}, err
	}

	relationships := make([]domain.Relationship, len(records))
	for i, record := range records {
		relationships[i] = MapRelationshipRecordToRelationship(record)
	}

	return relationships, nil
}

func (r *RelationshipRepository) UpdateRelationship(ctx context.Context, relationship domain.Relationship) (domain.Relationship, error) {
	params := queries.UpdateRelationshipParams{
		Id:         relationship.ID,
		FromToolId: relationship.FromToolID,
		ToToolId:   relationship.ToToolID,
		Kind:       relationship.Kind,
	}

	jsonMetadata, err := json.Marshal(relationship.Metadata)
	if err != nil {
		return domain.Relationship{}, err
	}
	params.Metadata = string(jsonMetadata)

	record, err := queries.New(r.db).UpdateRelationship(ctx, params)
	if err != nil {
		return domain.Relationship{}, err
	}

	return MapRelationshipRecordToRelationship(record), nil
}

func (r *RelationshipRepository) GetToolAlternatives(ctx context.Context, toolId int) ([]domain.Relationship, error) {
	params := queries.GetToolAlternativesParams{
		FromToolId: toolId,
		ToToolId:   toolId,
	}

	records, err := queries.New(r.db).GetToolAlternatives(ctx, params)
	if err != nil {
		return []domain.Relationship{}, err
	}

	relationships := make([]domain.Relationship, len(records))
	for i, record := range records {
		relationships[i] = MapRelationshipRecordToRelationship(record)
	}

	return relationships, nil
}

func (r *RelationshipRepository) GetRelationshipsAll(ctx context.Context, params GetRelationshipsAllParams) (PaginatedRelationshipsResult, error) {
	var createdAtString = ""
	if params.Cursor > 0 {
		createdAtString = time.Unix(params.Cursor, 0).UTC().Format(SQL_DATE_FORMAT)
	}

	limit := int(math.Min(float64(params.Limit), float64(MAX_LIMIT)))
	if limit == 0 {
		limit = DEFAULT_LIMIT
	}

	queryParams := queries.QueryParams{
		CreatedAt: createdAtString,
		Kind:      params.Kind,
		ToolId:    params.ToolId,
		Limit:     limit + 1,
	}
	relationships, totalCount, err := queries.GetPaginatedRelationships(ctx, r.db, queryParams)

	if err != nil {
		return PaginatedRelationshipsResult{}, err
	}

	result := make([]domain.Relationship, len(relationships))
	for i, record := range relationships {
		result[i] = MapRelationshipRecordToRelationship(record)
	}

	var nextCursor int64 = -1
	if len(result) > limit {
		nextCursor = result[len(result)-1].CreatedAt.Unix()
		result = result[:len(result)-1]
	}

	return PaginatedRelationshipsResult{
		Relationships: result,
		TotalCount:    totalCount,
		ItemsCount:    int64(len(result)),
		NextCursor:    nextCursor,
	}, nil
}
