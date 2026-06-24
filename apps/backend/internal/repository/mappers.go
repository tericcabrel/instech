package repository

import (
	"database/sql"
	"encoding/json"
	"tericcabrel/instech/db/queries"
	"tericcabrel/instech/internal/domain"
	"time"
)

const SQL_DATE_FORMAT = "2006-01-02 15:04:05"

func parseJSONArray[T any](jsonArray string) []T {
	var result []T
	err := json.Unmarshal([]byte(jsonArray), &result)
	if err != nil {
		return []T{}
	}

	return result
}

func parseJSONObject[T any](jsonObject string) T {
	var result T
	err := json.Unmarshal([]byte(jsonObject), &result)
	if err != nil {
		return *new(T)
	}

	return result
}

func optionalStringToNull(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: *value, Valid: true}
}

func nullStringToOptional(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	s := value.String

	return &s
}

func MapToolRecordToTool(tool queries.ToolRecord) (domain.Tool, error) {
	createdAt, err := time.Parse(SQL_DATE_FORMAT, tool.CreatedAt)
	if err != nil {
		return domain.Tool{}, err
	}
	updatedAt, err := time.Parse(SQL_DATE_FORMAT, tool.UpdatedAt)
	if err != nil {
		return domain.Tool{}, err
	}

	var mappedTool = domain.Tool{
		Id:          tool.Id,
		Name:        tool.Name,
		Slug:        tool.Slug,
		Category:    tool.Category,
		ReleaseYear: tool.ReleaseYear,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	if tool.SubType.Valid {
		mappedTool.SubType = tool.SubType.String
	}
	mappedTool.Prolang = nullStringToOptional(tool.Prolang)
	if tool.DevStatus.Valid {
		mappedTool.DevStatus = tool.DevStatus.String
	}
	mappedTool.Details = nullStringToOptional(tool.Details)
	mappedTool.Website = nullStringToOptional(tool.Website)
	mappedTool.Github = nullStringToOptional(tool.Github)

	mappedTool.UseCases = parseJSONArray[string](tool.UseCases)
	mappedTool.Tags = parseJSONArray[string](tool.Tags)

	return mappedTool, nil
}

func MapRelationshipRecordToRelationship(record queries.RelationshipRecord) domain.Relationship {
	createdAt, err := time.Parse(SQL_DATE_FORMAT, record.CreatedAt)
	if err != nil {
		return domain.Relationship{}
	}
	updatedAt, err := time.Parse(SQL_DATE_FORMAT, record.UpdatedAt)
	if err != nil {
		return domain.Relationship{}
	}
	var mappedRelationship = domain.Relationship{
		ID:         record.Id,
		FromToolID: record.FromToolId,
		ToToolID:   record.ToToolId,
		Kind:       record.Kind,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	if len(record.Metadata) > 0 {
		mappedRelationship.Metadata = parseJSONObject[domain.RelationshipMetadata](record.Metadata)
	}

	return mappedRelationship
}
