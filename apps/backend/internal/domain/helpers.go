package domain

func DedupeToolIdsFromRelationships(relationships []Relationship) []int {
	var seen = make(map[int]bool)
	for _, relationship := range relationships {
		seen[relationship.FromToolID] = true
		seen[relationship.ToToolID] = true
	}

	var toolIDs = make([]int, 0, len(seen))
	for toolID := range seen {
		toolIDs = append(toolIDs, toolID)
	}

	return toolIDs
}
