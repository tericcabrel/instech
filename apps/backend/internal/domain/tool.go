package domain

import (
	"slices"
	"strings"
	"time"
)

type Tool struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Category    string    `json:"category"`
	SubType     string    `json:"sub_type"`
	Prolang     string    `json:"prolang"`
	ReleaseYear int       `json:"release_year"`
	Devstatus   string    `json:"devstatus"`
	Details     string    `json:"details"`
	UseCases    []string  `json:"use_cases"`
	Tags        []string  `json:"tags"`
	Website     string    `json:"website"`
	Github      string    `json:"github"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateToolInput struct {
	Name        string
	Slug        string
	Category    string
	SubType     string
	Prolang     string
	ReleaseYear int
	Devstatus   string
	Details     string
	UseCases    []string
	Tags        []string
	Website     string
	Github      string
}

type UpdateToolInput struct {
	Name        string
	Slug        string
	Category    string
	SubType     string
	Prolang     string
	ReleaseYear int
	Devstatus   string
	Details     string
	UseCases    []string
	Tags        []string
	Website     string
	Github      string
}

var TOOL_CATEGORIES = []string{"language", "framework", "library"}
var TOOL_SUBTYPES = []string{"backend", "frontend", "fullstack", "mobile", "desktop", "game", "other"}
var TOOL_DEVSTATUSES = []string{"active", "deprecated"}

func IsCategoryValid(category string) bool {
	return slices.Contains(TOOL_CATEGORIES, category)
}
func IsSubTypeValid(subType string) bool {
	return slices.Contains(TOOL_SUBTYPES, subType)
}
func IsDevstatusValid(devstatus string) bool {
	return slices.Contains(TOOL_DEVSTATUSES, devstatus)
}

func areStringsEqual(a, b []string) bool {
	return slices.Equal(a, b)
}

func CreateTool(input CreateToolInput) (Tool, error) {
	tool := Tool{
		Name:        input.Name,
		Slug:        input.Slug,
		Category:    input.Category,
		SubType:     input.SubType,
		Prolang:     input.Prolang,
		ReleaseYear: input.ReleaseYear,
		Devstatus:   input.Devstatus,
		Details:     input.Details,
		UseCases:    input.UseCases,
		Tags:        input.Tags,
		Website:     input.Website,
		Github:      input.Github,
	}
	if !IsCategoryValid(input.Category) {
		return Tool{}, ErrInvalidToolCategory{Category: input.Category, Message: "The tool category is invalid. Valid categories are: " + strings.Join(TOOL_CATEGORIES, ", ")}
	}
	if !IsSubTypeValid(input.SubType) {
		return Tool{}, ErrInvalidToolSubType{SubType: input.SubType, Message: "The tool sub type is invalid. Valid sub types are: " + strings.Join(TOOL_SUBTYPES, ", ")}
	}
	if !IsDevstatusValid(input.Devstatus) {
		return Tool{}, ErrInvalidToolDevstatus{Devstatus: input.Devstatus, Message: "The tool dev status is invalid. Valid dev statuses are: " + strings.Join(TOOL_DEVSTATUSES, ", ")}
	}

	return tool, nil
}

func (tool *Tool) Update(input UpdateToolInput) error {
	if !IsCategoryValid(input.Category) {
		return ErrInvalidToolCategory{Category: input.Category, Message: "The tool category is invalid. Valid categories are: " + strings.Join(TOOL_CATEGORIES, ", ")}
	}
	if !IsSubTypeValid(input.SubType) {
		return ErrInvalidToolSubType{SubType: input.SubType, Message: "The tool sub type is invalid. Valid sub types are: " + strings.Join(TOOL_SUBTYPES, ", ")}
	}
	if !IsDevstatusValid(input.Devstatus) {
		return ErrInvalidToolDevstatus{Devstatus: input.Devstatus, Message: "The tool dev status is invalid. Valid dev statuses are: " + strings.Join(TOOL_DEVSTATUSES, ", ")}
	}

	if tool.Name != input.Name {
		tool.Name = input.Name
	}
	if tool.Slug != input.Slug {
		tool.Slug = input.Slug
	}
	if tool.Category != input.Category {
		tool.Category = input.Category
	}
	if tool.SubType != input.SubType {
		tool.SubType = input.SubType
	}
	if tool.Prolang != input.Prolang {
		tool.Prolang = input.Prolang
	}
	if tool.ReleaseYear != input.ReleaseYear {
		tool.ReleaseYear = input.ReleaseYear
	}
	if tool.Devstatus != input.Devstatus {
		tool.Devstatus = input.Devstatus
	}
	if tool.Details != input.Details {
		tool.Details = input.Details
	}
	if !areStringsEqual(tool.UseCases, input.UseCases) {
		tool.UseCases = input.UseCases
	}
	if !areStringsEqual(tool.Tags, input.Tags) {
		tool.Tags = input.Tags
	}
	if tool.Website != input.Website {
		tool.Website = input.Website
	}
	if tool.Github != input.Github {
		tool.Github = input.Github
	}

	return nil
}
