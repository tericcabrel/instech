package domain_test

import (
	"slices"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/testutil"
	"testing"
)

func TestTool(t *testing.T) {
	t.Run("Create tool with invalid category", func(t *testing.T) {
		_, err := domain.CreateTool(domain.CreateToolInput{
			Name:     "Test Tool",
			Category: "invalid",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidToolCategory); !ok {
			t.Errorf("Expected ErrInvalidToolCategory, got %v", err)
		}
	})
	t.Run("Create tool with invalid sub type", func(t *testing.T) {
		_, err := domain.CreateTool(domain.CreateToolInput{
			Name:     "Test Tool",
			Category: "language",
			SubType:  "invalid",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidToolSubType); !ok {
			t.Errorf("Expected ErrInvalidToolSubType, got %v", err)
		}
	})

	t.Run("Create tool with invalid dev status", func(t *testing.T) {
		_, err := domain.CreateTool(domain.CreateToolInput{
			Name:      "Test Tool",
			Category:  "language",
			SubType:   "backend",
			DevStatus: "invalid",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidToolDevStatus); !ok {
			t.Errorf("Expected ErrInvalidToolDevStatus, got %v", err)
		}
	})

	t.Run("Create tool with invalid fields", func(t *testing.T) {
		_, err := domain.CreateTool(domain.CreateToolInput{
			Name:        "",
			Category:    "language",
			SubType:     "backend",
			DevStatus:   "active",
			Details:     new(""),
			UseCases:    []string{},
			Tags:        []string{},
			Website:     new("not-a-valid-url"),
			Github:      new("not-a-valid-url"),
			ReleaseYear: 1929,
			Prolang:     new(""),
			Slug:        "",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidField); ok {
			fields := []string{"Name", "Slug", "ReleaseYear", "Prolang", "Details", "Website", "Github"}
			for _, field := range fields {
				if _, exist := e.Fields[field]; !exist {
					t.Errorf("Expected the field \"%s\" to be present", field)
				}
			}
		}
	})

	t.Run("Create tool with valid input", func(t *testing.T) {
		tool, err := domain.CreateTool(domain.CreateToolInput{
			Name:        "Test Tool",
			Category:    "language",
			SubType:     "backend",
			DevStatus:   "active",
			Details:     new("Test Details"),
			UseCases:    []string{"Test Use Case"},
			Tags:        []string{"Test Tag"},
			Website:     new("https://test.com"),
			Github:      new("https://github.com/test"),
			ReleaseYear: 2020,
			Prolang:     new("Go"),
			Slug:        "golang",
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if tool.Name != "Test Tool" {
			t.Errorf("Expected name to be 'Test Tool', got %s", tool.Name)
		}
		if tool.Category != "language" {
			t.Errorf("Expected category to be 'language', got %s", tool.Category)
		}
		if tool.SubType != "backend" {
			t.Errorf("Expected sub type to be 'backend', got %s", tool.SubType)
		}
		if tool.DevStatus != "active" {
			t.Errorf("Expected dev status to be 'active', got %s", tool.DevStatus)
		}
		if tool.Details == nil || *tool.Details != "Test Details" {
			t.Errorf("Expected details to be 'Test Details', got %v", tool.Details)
		}
		if !slices.Equal(tool.UseCases, []string{"Test Use Case"}) {
			t.Errorf("Expected use cases to be ['Test Use Case'], got %v", tool.UseCases)
		}
		if !slices.Equal(tool.Tags, []string{"Test Tag"}) {
			t.Errorf("Expected tags to be ['Test Tag'], got %v", tool.Tags)
		}
		if tool.Website == nil || *tool.Website != "https://test.com" {
			t.Errorf("Expected website to be 'https://test.com', got %v", tool.Website)
		}
		if tool.Github == nil || *tool.Github != "https://github.com/test" {
			t.Errorf("Expected github to be 'https://github.com/test', got %v", tool.Github)
		}
		if tool.ReleaseYear != 2020 {
			t.Errorf("Expected release year to be 2020, got %d", tool.ReleaseYear)
		}
		if tool.Prolang == nil || *tool.Prolang != "Go" {
			t.Errorf("Expected pro lang to be 'Go', got %v", tool.Prolang)
		}
		if tool.Slug != "golang" {
			t.Errorf("Expected slug to be 'golang', got %s", tool.Slug)
		}
	})

	t.Run("Create tool allows omitted optional string fields", func(t *testing.T) {
		tool, err := domain.CreateTool(domain.CreateToolInput{
			Name:        "React",
			Slug:        "react",
			Category:    "framework",
			SubType:     "frontend",
			DevStatus:   "active",
			ReleaseYear: 2013,
		})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if tool.Prolang != nil {
			t.Errorf("Expected prolang to be nil, got %v", tool.Prolang)
		}
		if tool.Details != nil {
			t.Errorf("Expected details to be nil, got %v", tool.Details)
		}
		if tool.Website != nil {
			t.Errorf("Expected website to be nil, got %v", tool.Website)
		}
		if tool.Github != nil {
			t.Errorf("Expected github to be nil, got %v", tool.Github)
		}
	})

	t.Run("Create tool rejects empty optional string fields", func(t *testing.T) {
		_, err := domain.CreateTool(domain.CreateToolInput{
			Name:        "React",
			Slug:        "react",
			Category:    "framework",
			SubType:     "frontend",
			DevStatus:   "active",
			ReleaseYear: 2013,
			Details:     new(""),
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if e, ok := err.(domain.ErrInvalidField); !ok {
			t.Fatalf("Expected ErrInvalidField, got %v", err)
		} else if _, exists := e.Fields["Details"]; !exists {
			t.Fatalf("Expected Details field error, got %v", e.Fields)
		}
	})

	t.Run("Update tool with invalid category", func(t *testing.T) {
		tool := testutil.CreateTestTool()

		err := tool.Update(domain.UpdateToolInput{
			Category: "invalid",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidToolCategory); !ok {
			t.Errorf("Expected ErrInvalidToolCategory, got %v", err)
		}
	})

	t.Run("Update tool with invalid sub type", func(t *testing.T) {
		tool := testutil.CreateTestTool()

		err := tool.Update(domain.UpdateToolInput{
			Category: "language",
			SubType:  "invalid",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidToolSubType); !ok {
			t.Errorf("Expected ErrInvalidToolSubType, got %v", err)
		}
	})

	t.Run("Update tool with invalid dev status", func(t *testing.T) {
		tool := testutil.CreateTestTool()

		err := tool.Update(domain.UpdateToolInput{
			Category:  "language",
			SubType:   "backend",
			DevStatus: "invalid",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidToolDevStatus); !ok {
			t.Errorf("Expected ErrInvalidToolDevStatus, got %v", err)
		}
	})

	t.Run("Update tool with invalid fields", func(t *testing.T) {
		tool := testutil.CreateTestTool()

		err := tool.Update(domain.UpdateToolInput{
			Name:        "",
			Category:    "language",
			SubType:     "backend",
			DevStatus:   "active",
			Details:     new(""),
			UseCases:    []string{},
			Tags:        []string{},
			Website:     new("not-a-valid-url"),
			Github:      new("not-a-valid-url"),
			ReleaseYear: 1929,
			Prolang:     new(""),
			Slug:        "",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidField); ok {
			fields := []string{"Name", "Slug", "ReleaseYear", "Prolang", "Details", "Website", "Github"}
			for _, field := range fields {
				if _, exist := e.Fields[field]; !exist {
					t.Errorf("Expected the field \"%s\" to be present", field)
				}
			}
		}
	})

	t.Run("Update tool with valid input", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		updateInput := domain.UpdateToolInput{
			Name:        "Mootools",
			Category:    "framework",
			SubType:     "frontend",
			DevStatus:   "deprecated",
			Details:     new("Mootools is a framework"),
			UseCases:    []string{"SPA", "SEO", "API", "Frontend"},
			Tags:        []string{"component-based", "declarative", "functional", "object-oriented"},
			Website:     new("https://mootools.net"),
			Github:      new("https://github.com/mootools/mootools"),
			ReleaseYear: 2006,
			Prolang:     new("javascript"),
			Slug:        "mootools",
		}

		err := tool.Update(updateInput)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if tool.Name != "Mootools" {
			t.Errorf("Expected name to be 'Mootools', got %s", tool.Name)
		}
		if tool.Category != "framework" {
			t.Errorf("Expected category to be 'framework', got %s", tool.Category)
		}
		if tool.SubType != "frontend" {
			t.Errorf("Expected sub type to be 'frontend', got %s", tool.SubType)
		}
		if tool.DevStatus != "deprecated" {
			t.Errorf("Expected dev status to be 'deprecated', got %s", tool.DevStatus)
		}
		if tool.Details == nil || *tool.Details != "Mootools is a framework" {
			t.Errorf("Expected details to be 'Mootools is a framework', got %v", tool.Details)
		}
		if !slices.Equal(tool.UseCases, []string{"SPA", "SEO", "API", "Frontend"}) {
			t.Errorf("Expected use cases to be ['SPA', 'SEO', 'API', 'Frontend'], got %v", tool.UseCases)
		}
		if !slices.Equal(tool.Tags, []string{"component-based", "declarative", "functional", "object-oriented"}) {
			t.Errorf("Expected tags to be ['component-based', 'declarative', 'functional', 'object-oriented'], got %v", tool.Tags)
		}
		if tool.Website == nil || *tool.Website != "https://mootools.net" {
			t.Errorf("Expected website to be 'https://mootools.net', got %v", tool.Website)
		}
		if tool.Github == nil || *tool.Github != "https://github.com/mootools/mootools" {
			t.Errorf("Expected github to be 'https://github.com/mootools/mootools', got %v", tool.Github)
		}
		if tool.ReleaseYear != 2006 {
			t.Errorf("Expected release year to be 2006, got %d", tool.ReleaseYear)
		}
		if tool.Prolang == nil || *tool.Prolang != "javascript" {
			t.Errorf("Expected pro lang to be 'javascript', got %v", tool.Prolang)
		}
		if tool.Slug != "mootools" {
			t.Errorf("Expected slug to be 'mootools', got %s", tool.Slug)
		}
	})

	t.Run("Update tool with valid input and reset optional fields", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		updateInput := domain.UpdateToolInput{
			Name:        "Mootools",
			Category:    "framework",
			SubType:     "frontend",
			DevStatus:   "deprecated",
			Details:     nil,
			UseCases:    []string{"SPA", "SEO", "API", "Frontend"},
			Tags:        []string{"component-based", "declarative", "functional", "object-oriented"},
			Website:     nil,
			Github:      nil,
			ReleaseYear: 2006,
			Prolang:     new("javascript"),
			Slug:        "mootools",
		}

		err := tool.Update(updateInput)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if tool.Details != nil {
			t.Errorf("Expected details to be nil, got %v", tool.Details)
		}
		if tool.Website != nil {
			t.Errorf("Expected website to be nil, got %v", tool.Website)
		}
		if tool.Github != nil {
			t.Errorf("Expected github to be nil, got %v", tool.Github)
		}
	})
}
