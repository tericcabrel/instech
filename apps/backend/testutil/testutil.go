package testutil

import "tericcabrel/instech/internal/domain"

func CreateTestTool() domain.Tool {
	const ReleaseYear = 2020
	tool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Golang",
		Category:    "language",
		SubType:     "backend",
		DevStatus:   "active",
		Details:     new("Test Details"),
		UseCases:    []string{"Test Use Case"},
		Tags:        []string{"Test Tag"},
		Website:     new("https://test.com"),
		Github:      new("https://github.com/test"),
		ReleaseYear: ReleaseYear,
		Prolang:     new("Go"),
		Slug:        "golang",
	})

	if err != nil {
		panic(err)
	}

	tool.Id = 1

	return tool
}

func CreateTestDynamicTool(id int, input domain.CreateToolInput) domain.Tool {
	tool, err := domain.CreateTool(input)

	if err != nil {
		panic(err)
	}
	tool.Id = id

	return tool
}

func CreateTestRelationship() domain.Relationship {
	const fromToolID = 1
	const toToolID = 2
	relationship, err := domain.CreateRelationship(domain.CreateRelationshipInput{
		FromToolID: fromToolID,
		ToToolID:   toToolID,
		Kind:       "built_on",
		Reason:     "This is a test relationship",
	})

	if err != nil {
		panic(err)
	}

	relationship.ID = 1

	return relationship
}

func CreateTestDynamicRelationship(id int, input domain.CreateRelationshipInput) domain.Relationship {
	relationship, err := domain.CreateRelationship(input)

	if err != nil {
		panic(err)
	}
	relationship.ID = id

	return relationship
}
