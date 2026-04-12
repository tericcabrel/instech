package domain

import "time"

type RelationshipMetadata struct {
	Reason string `json:"reason"`
}

type Relationship struct {
	Id         int                  `json:"id"`
	FromToolId int                  `json:"from_tool_id"`
	ToToolId   int                  `json:"to_tool_id"`
	Kind       string               `json:"kind"`
	Metadata   RelationshipMetadata `json:"metadata"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

var RELATIONSHIP_KINDS = []string{"built_on", "inspired_by", "alternative_to", "replaced_by", "used_with"}
