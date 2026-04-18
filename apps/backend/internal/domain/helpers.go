package domain

func DedupeToolIdsFromRelationships(relationships []Relationship) []int {
	var seen map[int]bool = make(map[int]bool)
	for _, relationship := range relationships {
		seen[relationship.FromToolId] = true
		seen[relationship.ToToolId] = true
	}

	var toolIds []int = make([]int, 0)
	for toolId := range seen {
		toolIds = append(toolIds, toolId)
	}

	return toolIds
}
