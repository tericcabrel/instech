package testutil

import "tericcabrel/instech/internal/domain"

func CreateTestTool() domain.Tool {
	tool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Golang",
		Category:    "language",
		SubType:     "backend",
		Devstatus:   "active",
		Details:     "Test Details",
		UseCases:    []string{"Test Use Case"},
		Tags:        []string{"Test Tag"},
		Website:     "https://test.com",
		Github:      "https://github.com/test",
		ReleaseYear: 2020,
		Prolang:     "Go",
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
	relationship, err := domain.CreateRelationship(domain.CreateRelationshipInput{
		FromToolId: 1,
		ToToolId:   2,
		Kind:       "built_on",
		Reason:     "This is a test relationship",
	})

	if err != nil {
		panic(err)
	}

	relationship.Id = 1

	return relationship
}

func CreateTestDynamicRelationship(id int, input domain.CreateRelationshipInput) domain.Relationship {
	relationship, err := domain.CreateRelationship(input)

	if err != nil {
		panic(err)
	}
	relationship.Id = id

	return relationship
}
