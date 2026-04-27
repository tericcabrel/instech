package domain

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
)

type Tool struct {
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Website     string    `json:"website"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Category    string    `json:"category"`
	SubType     string    `json:"sub_type"`
	Prolang     string    `json:"prolang"`
	Devstatus   string    `json:"devstatus"`
	Details     string    `json:"details"`
	Github      string    `json:"github"`
	UseCases    []string  `json:"use_cases"`
	Tags        []string  `json:"tags"`
	Id          int       `json:"id"`
	ReleaseYear int       `json:"release_year"`
}

type CreateToolInput struct {
	Name        string
	Slug        string
	Category    string
	SubType     string
	Prolang     string
	Devstatus   string
	Details     string
	Website     string
	Github      string
	UseCases    []string
	Tags        []string
	ReleaseYear int
}

type UpdateToolInput struct {
	Name        string
	Slug        string
	Category    string
	SubType     string
	Prolang     string
	Devstatus   string
	Details     string
	Website     string
	Github      string
	UseCases    []string
	Tags        []string
	ReleaseYear int
}

var ToolCategories = []string{"language", "framework", "library"}
var ToolSubtypes = []string{"backend", "frontend", "fullstack", "mobile", "desktop", "game", "other"}
var ToolDevStatuses = []string{"active", "deprecated"}

const MIN_RELEASE_YEAR = 1940

var MAX_RELEASE_YEAR = time.Now().Year()

const ERROR_NAME_REQUIRED = "The tool name is required"
const ERROR_SLUG_REQUIRED = "The tool slug is required"
const ERROR_PROLANG_REQUIRED = "The tool programming language is required"
const ERROR_WEBSITE_INVALID = "The tool website is invalid. Valid websites must be a valid URL"
const ERROR_GITHUB_INVALID = "The tool github is invalid. Valid github must be a valid URL"
const ERROR_USE_CASES_INVALID = "The tool use cases are invalid. Valid use cases must be an array of strings"
const ERROR_TAGS_INVALID = "The tool tags are invalid. Valid tags must be an array of strings"

var ERROR_RELEASE_YEAR_INVALID = fmt.Sprintf("The tool release year is invalid. Valid release years are between %d and %d", MIN_RELEASE_YEAR, MAX_RELEASE_YEAR)
var ERROR_CATEGORY_INVALID = fmt.Sprintf("The tool category is invalid. Valid categories are: %s", strings.Join(ToolCategories, ", "))
var ERROR_SUBTYPE_INVALID = fmt.Sprintf("The tool sub type is invalid. Valid sub types are: %s", strings.Join(ToolSubtypes, ", "))
var ERROR_DEVSTATUS_INVALID = fmt.Sprintf("The tool dev status is invalid. Valid dev statuses are: %s", strings.Join(ToolDevStatuses, ", "))

func IsCategoryValid(category string) bool {
	return slices.Contains(ToolCategories, category)
}
func IsSubTypeValid(subType string) bool {
	return slices.Contains(ToolSubtypes, subType)
}
func IsDevstatusValid(devstatus string) bool {
	return slices.Contains(ToolDevStatuses, devstatus)
}

func areStringsEqual(a, b []string) bool {
	return slices.Equal(a, b)
}

func isValidURL(url string) bool {
	const REGEX_URL = `^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`
	regex := regexp.MustCompile(REGEX_URL)
	return regex.MatchString(url)
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
		return Tool{}, ErrInvalidToolCategory{Category: input.Category, Message: ERROR_CATEGORY_INVALID}
	}
	if !IsSubTypeValid(input.SubType) {
		return Tool{}, ErrInvalidToolSubType{SubType: input.SubType, Message: ERROR_SUBTYPE_INVALID}
	}
	if !IsDevstatusValid(input.Devstatus) {
		return Tool{}, ErrInvalidToolDevstatus{Devstatus: input.Devstatus, Message: ERROR_DEVSTATUS_INVALID}
	}

	var errors = make(map[string]string)
	if input.Name == "" {
		errors["Name"] = ERROR_NAME_REQUIRED
	}
	if input.Slug == "" {
		errors["Slug"] = ERROR_SLUG_REQUIRED
	}

	if input.ReleaseYear < 1940 || input.ReleaseYear > time.Now().Year() {
		errors["ReleaseYear"] = ERROR_RELEASE_YEAR_INVALID
	}

	if input.Category == "language" && input.Prolang == "" {
		errors["Prolang"] = ERROR_PROLANG_REQUIRED
	}

	if input.Website != "" && !isValidURL(input.Website) {
		errors["Website"] = ERROR_WEBSITE_INVALID
	}
	if input.Github != "" && !isValidURL(input.Github) {
		errors["Github"] = ERROR_GITHUB_INVALID
	}

	if len(errors) > 0 {
		return Tool{}, ErrInvalidField{Fields: errors}
	}

	return tool, nil
}

func (tool *Tool) Update(input UpdateToolInput) error {
	if !IsCategoryValid(input.Category) {
		return ErrInvalidToolCategory{Category: input.Category, Message: ERROR_CATEGORY_INVALID}
	}
	if !IsSubTypeValid(input.SubType) {
		return ErrInvalidToolSubType{SubType: input.SubType, Message: ERROR_SUBTYPE_INVALID}
	}
	if !IsDevstatusValid(input.Devstatus) {
		return ErrInvalidToolDevstatus{Devstatus: input.Devstatus, Message: ERROR_DEVSTATUS_INVALID}
	}

	var errors = make(map[string]string)

	if tool.Name != input.Name {
		if input.Name == "" {
			errors["Name"] = ERROR_NAME_REQUIRED
		} else {
			tool.Name = input.Name
		}
	}
	if tool.Slug != input.Slug {
		if input.Slug == "" {
			errors["Slug"] = ERROR_SLUG_REQUIRED
		} else {
			tool.Slug = input.Slug
		}
	}
	if tool.Category != input.Category {
		tool.Category = input.Category
	}
	if tool.SubType != input.SubType {
		tool.SubType = input.SubType
	}
	if tool.Prolang != input.Prolang {
		if input.Category == "language" && input.Prolang == "" {
			errors["Prolang"] = ERROR_PROLANG_REQUIRED
		} else {
			tool.Prolang = input.Prolang
		}
	}
	if tool.ReleaseYear != input.ReleaseYear {
		if input.ReleaseYear < 1940 || input.ReleaseYear > time.Now().Year() {
			errors["ReleaseYear"] = ERROR_RELEASE_YEAR_INVALID
		} else {
			tool.ReleaseYear = input.ReleaseYear
		}
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
		if input.Website != "" && !isValidURL(input.Website) {
			errors["Website"] = ERROR_WEBSITE_INVALID
		} else {
			tool.Website = input.Website
		}
	}
	if tool.Github != input.Github {
		if input.Github != "" && !isValidURL(input.Github) {
			errors["Github"] = ERROR_GITHUB_INVALID
		} else {
			tool.Github = input.Github
		}
	}

	if len(errors) > 0 {
		return ErrInvalidField{Fields: errors}
	}

	return nil
}
