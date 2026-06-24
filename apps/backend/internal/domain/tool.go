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
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Category    string    `json:"category"`
	SubType     string    `json:"sub_type"`
	Prolang     *string   `json:"prolang,omitempty"`
	DevStatus   string    `json:"dev_status"`
	Details     *string   `json:"details,omitempty"`
	Website     *string   `json:"website,omitempty"`
	Github      *string   `json:"github,omitempty"`
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
	Prolang     *string
	DevStatus   string
	Details     *string
	Website     *string
	Github      *string
	UseCases    []string
	Tags        []string
	ReleaseYear int
}

type UpdateToolInput struct {
	Name        string
	Slug        string
	Category    string
	SubType     string
	Prolang     *string
	DevStatus   string
	Details     *string
	Website     *string
	Github      *string
	UseCases    []string
	Tags        []string
	ReleaseYear int
}

var ToolCategories = []string{"language", "framework", "library"}
var ToolSubtypes = []string{"backend", "frontend", "fullstack", "mobile", "desktop", "game", "other"}
var ToolDevStatuses = []string{"active", "deprecated"}

const MIN_RELEASE_YEAR = 1940

var MAX_RELEASE_YEAR = time.Now().Year()

const ERROR_NAME_REQUIRED = "The name is required"
const ERROR_SLUG_REQUIRED = "The slug is required"
const ERROR_PROLANG_REQUIRED = "The programming language is required"
const ERROR_PROLANG_EMPTY = "The programming language cannot be empty"
const ERROR_DETAILS_EMPTY = "The details cannot be empty"
const ERROR_WEBSITE_EMPTY = "The website cannot be empty"
const ERROR_WEBSITE_INVALID = "The website is invalid. Valid websites must be a valid URL"
const ERROR_GITHUB_EMPTY = "The github cannot be empty"
const ERROR_GITHUB_INVALID = "The github is invalid. Valid github must be a valid URL"
const ERROR_USE_CASES_INVALID = "The use cases are invalid. Valid use cases must be an array of strings"
const ERROR_TAGS_INVALID = "The tags are invalid. Valid tags must be an array of strings"

var ERROR_RELEASE_YEAR_INVALID = fmt.Sprintf("The release year is invalid. Valid release years are between %d and %d", MIN_RELEASE_YEAR, MAX_RELEASE_YEAR)
var ERROR_CATEGORY_INVALID = fmt.Sprintf("The category is invalid. Valid categories are: %s", strings.Join(ToolCategories, ", "))
var ERROR_SUBTYPE_INVALID = fmt.Sprintf("The sub type is invalid. Valid sub types are: %s", strings.Join(ToolSubtypes, ", "))
var ERROR_DEVSTATUS_INVALID = fmt.Sprintf("The dev status is invalid. Valid dev statuses are: %s", strings.Join(ToolDevStatuses, ", "))

func IsCategoryValid(category string) bool {
	return slices.Contains(ToolCategories, category)
}
func IsSubTypeValid(subType string) bool {
	return slices.Contains(ToolSubtypes, subType)
}
func IsDevStatusValid(devStatus string) bool {
	return slices.Contains(ToolDevStatuses, devStatus)
}

func areStringsEqual(a, b []string) bool {
	return slices.Equal(a, b)
}

func optionalStringEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

func validateOptionalStringNotEmpty(fieldKey string, value *string, errors map[string]string, emptyMessage string) {
	if value != nil && *value == "" {
		errors[fieldKey] = emptyMessage
	}
}

func validateProlang(category string, prolang *string, errors map[string]string) {
	validateOptionalStringNotEmpty("Prolang", prolang, errors, ERROR_PROLANG_EMPTY)
	if category == "language" && prolang == nil {
		errors["Prolang"] = ERROR_PROLANG_REQUIRED
	}
}

func validateWebsite(website *string, errors map[string]string) {
	validateOptionalStringNotEmpty("Website", website, errors, ERROR_WEBSITE_EMPTY)
	if website != nil && *website != "" && !isValidURL(*website) {
		errors["Website"] = ERROR_WEBSITE_INVALID
	}
}

func validateGithub(github *string, errors map[string]string) {
	validateOptionalStringNotEmpty("Github", github, errors, ERROR_GITHUB_EMPTY)
	if github != nil && *github != "" && !isValidURL(*github) {
		errors["Github"] = ERROR_GITHUB_INVALID
	}
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
		DevStatus:   input.DevStatus,
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
	if !IsDevStatusValid(input.DevStatus) {
		return Tool{}, ErrInvalidToolDevStatus{DevStatus: input.DevStatus, Message: ERROR_DEVSTATUS_INVALID}
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

	validateProlang(input.Category, input.Prolang, errors)
	validateOptionalStringNotEmpty("Details", input.Details, errors, ERROR_DETAILS_EMPTY)
	validateWebsite(input.Website, errors)
	validateGithub(input.Github, errors)

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
	if !IsDevStatusValid(input.DevStatus) {
		return ErrInvalidToolDevStatus{DevStatus: input.DevStatus, Message: ERROR_DEVSTATUS_INVALID}
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
	if !optionalStringEqual(tool.Prolang, input.Prolang) {
		validateProlang(input.Category, input.Prolang, errors)
		if _, hasProlangError := errors["Prolang"]; !hasProlangError {
			tool.Prolang = input.Prolang
		}
	} else if input.Category == "language" && input.Prolang == nil {
		errors["Prolang"] = ERROR_PROLANG_REQUIRED
	}
	if tool.ReleaseYear != input.ReleaseYear {
		if input.ReleaseYear < 1940 || input.ReleaseYear > time.Now().Year() {
			errors["ReleaseYear"] = ERROR_RELEASE_YEAR_INVALID
		} else {
			tool.ReleaseYear = input.ReleaseYear
		}
	}
	if tool.DevStatus != input.DevStatus {
		tool.DevStatus = input.DevStatus
	}
	if !optionalStringEqual(tool.Details, input.Details) {
		validateOptionalStringNotEmpty("Details", input.Details, errors, ERROR_DETAILS_EMPTY)
		if _, hasDetailsError := errors["Details"]; !hasDetailsError {
			tool.Details = input.Details
		}
	}
	if !areStringsEqual(tool.UseCases, input.UseCases) {
		tool.UseCases = input.UseCases
	}
	if !areStringsEqual(tool.Tags, input.Tags) {
		tool.Tags = input.Tags
	}
	if !optionalStringEqual(tool.Website, input.Website) {
		validateWebsite(input.Website, errors)
		if _, hasWebsiteError := errors["Website"]; !hasWebsiteError {
			tool.Website = input.Website
		}
	}
	if !optionalStringEqual(tool.Github, input.Github) {
		validateGithub(input.Github, errors)
		if _, hasGithubError := errors["Github"]; !hasGithubError {
			tool.Github = input.Github
		}
	}

	if len(errors) > 0 {
		return ErrInvalidField{Fields: errors}
	}

	return nil
}
