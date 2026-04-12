package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"tericcabrel/instech/db/queries"
	"tericcabrel/instech/internal/domain"
)

type RelationshipRepositoryInterface interface {
	CreateRelationship(ctx context.Context, relationship domain.Relationship) (domain.Relationship, error)
	DeleteRelationship(ctx context.Context, id int) error
	GetRelationshipsByToolID(ctx context.Context, toolID int) ([]domain.Relationship, error)
	UpdateRelationship(ctx context.Context, relationship domain.Relationship) (domain.Relationship, error)
}

type RelationshipRepository struct {
	db *sql.DB
}

func NewRelationshipRepository(db *sql.DB) *RelationshipRepository {
	return &RelationshipRepository{db: db}
}

func (r *RelationshipRepository) CreateRelationship(ctx context.Context, relationship domain.Relationship) (domain.Relationship, error) {
	params := queries.CreateRelationshipParams{
		FromToolId: relationship.FromToolId,
		ToToolId:   relationship.ToToolId,
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
		Id:         relationship.Id,
		FromToolId: relationship.FromToolId,
		ToToolId:   relationship.ToToolId,
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
