package queries

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"
)

type QueryParams struct {
	CreatedAt string
	Kind      string
	ToolId    int
	Limit     int
}

type BuildQueryAndArgsBuilder struct {
	PaginationQuery     string
	CountQuery          string
	PaginationQueryArgs []any
	CountQueryArgs      []any
}

func buildQueryAndArgs(params QueryParams) BuildQueryAndArgsBuilder {
	paginationQuery := "SELECT * FROM relationships "
	countQuery := "SELECT COUNT(*) FROM relationships "
	var conditionsForPagination []string = make([]string, 0)
	var conditionsForCount []string = make([]string, 0)
	var paginationQueryArgs []any = make([]any, 0)
	var countQueryArgs []any = make([]any, 0)

	if params.CreatedAt != "" {
		conditionsForPagination = append(conditionsForPagination, "created_at <= ?")
		paginationQueryArgs = append(paginationQueryArgs, params.CreatedAt)
	}
	if params.Kind != "" {
		conditionsForPagination = append(conditionsForPagination, "kind = ?")
		paginationQueryArgs = append(paginationQueryArgs, params.Kind)
		conditionsForCount = append(conditionsForCount, "kind = ?")
		countQueryArgs = append(countQueryArgs, params.Kind)
	}
	if params.ToolId != 0 {
		conditionsForPagination = append(conditionsForPagination, "(from_tool_id = ? OR to_tool_id = ?)")
		paginationQueryArgs = append(paginationQueryArgs, params.ToolId)
		paginationQueryArgs = append(paginationQueryArgs, params.ToolId)
		conditionsForCount = append(conditionsForCount, "(from_tool_id = ? OR to_tool_id = ?)")
		countQueryArgs = append(countQueryArgs, params.ToolId)
		countQueryArgs = append(countQueryArgs, params.ToolId)
	}
	if len(conditionsForPagination) > 0 {
		paginationQuery += " WHERE " + strings.Join(conditionsForPagination, " AND ")
	}
	if len(conditionsForCount) > 0 {
		countQuery += " WHERE " + strings.Join(conditionsForCount, " AND ")
	}

	paginationQuery += " ORDER BY created_at DESC LIMIT ?"

	paginationQueryArgs = append(paginationQueryArgs, params.Limit)

	return BuildQueryAndArgsBuilder{
		PaginationQuery:     paginationQuery,
		CountQuery:          countQuery,
		PaginationQueryArgs: paginationQueryArgs,
		CountQueryArgs:      countQueryArgs,
	}
}

func GetPaginatedRelationships(ctx context.Context, db *sql.DB, params QueryParams) ([]RelationshipRecord, int64, error) {
	builder := buildQueryAndArgs(params)
	rows, err := db.QueryContext(ctx, builder.PaginationQuery, builder.PaginationQueryArgs...)

	if err != nil {
		return nil, 0, errors.Wrap(err, "GetPaginatedRelationships")
	}

	defer rows.Close()

	items := []RelationshipRecord{}
	for rows.Next() {
		var i RelationshipRecord
		if err := rows.Scan(&i.Id, &i.FromToolId, &i.ToToolId, &i.Kind, &i.Metadata, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, 0, errors.Wrap(err, "GetPaginatedRelationshipsRowsScan")
		}
		items = append(items, i)
	}

	row := db.QueryRowContext(ctx, builder.CountQuery, builder.CountQueryArgs...)
	var totalCount int64
	if err := row.Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "GetPaginatedRelationshipsTotalCount")
	}

	return items, totalCount, nil
}
