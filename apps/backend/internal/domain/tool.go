package domain

import "time"

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

var TOOL_CATEGORIES = []string{"language", "framework", "library"}
var TOOL_SUBTYPES = []string{"backend", "frontend", "fullstack", "mobile", "desktop", "game", "other"}
var TOOL_DEVSTATUSES = []string{"active", "deprecated"}
