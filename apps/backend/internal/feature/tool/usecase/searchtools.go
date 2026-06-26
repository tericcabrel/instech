package usecase

import (
	"context"
	"strings"
	"tericcabrel/instech/internal/repository"
)

type SearchToolsUseCase struct {
	ToolRepository repository.ToolRepositoryInterface
}

type SearchToolsResult struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Category string `json:"category"`
	SubType  string `json:"subType"`
	Id       int    `json:"id"`
}

func (uc *SearchToolsUseCase) Execute(keyword string) ([]SearchToolsResult, error) {
	normalizedKeyword := strings.TrimSpace(keyword)
	if normalizedKeyword == "" {
		return []SearchToolsResult{}, nil
	}

	results, err := uc.ToolRepository.SearchTools(context.Background(), normalizedKeyword)
	if err != nil {
		return []SearchToolsResult{}, err
	}

	var searchResults = make([]SearchToolsResult, 0, len(results))
	for _, result := range results {
		searchResults = append(searchResults, SearchToolsResult{
			Id:       result.Id,
			Slug:     result.Slug,
			Name:     result.Name,
			Category: result.Category,
			SubType:  result.SubType,
		})
	}

	return searchResults, nil
}
