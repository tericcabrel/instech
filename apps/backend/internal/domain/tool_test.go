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
			Devstatus: "invalid",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidToolDevstatus); !ok {
			t.Errorf("Expected ErrInvalidToolDevstatus, got %v", err)
		}
	})

	t.Run("Create tool with invalid fields", func(t *testing.T) {
		_, err := domain.CreateTool(domain.CreateToolInput{
			Name:        "",
			Category:    "language",
			SubType:     "backend",
			Devstatus:   "active",
			Details:     "",
			UseCases:    []string{},
			Tags:        []string{},
			Website:     "not-a-valid-url",
			Github:      "not-a-valid-url",
			ReleaseYear: 1929,
			Prolang:     "",
			Slug:        "",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidField); ok {
			if _, exist := e.Fields["Name"]; !exist {
				t.Errorf("Expected the field \"Name\" to be present")
			}
			if _, exist := e.Fields["Slug"]; !exist {
				t.Errorf("Expected the field \"Slug\" to be present")
			}
			if _, exist := e.Fields["ReleaseYear"]; !exist {
				t.Errorf("Expected the field \"ReleaseYear\" to be present")
			}
			if _, exist := e.Fields["Prolang"]; !exist {
				t.Errorf("Expected the field \"Prolang\" to be present")
			}
			if _, exist := e.Fields["Website"]; !exist {
				t.Errorf("Expected the field \"Website\" to be present")
			}
			if _, exist := e.Fields["Github"]; !exist {
				t.Errorf("Expected the field \"Github\" to be present")
			}
		}
	})

	t.Run("Create tool with valid input", func(t *testing.T) {
		tool, err := domain.CreateTool(domain.CreateToolInput{
			Name:        "Test Tool",
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
		if tool.Devstatus != "active" {
			t.Errorf("Expected dev status to be 'active', got %s", tool.Devstatus)
		}
		if tool.Details != "Test Details" {
			t.Errorf("Expected details to be 'Test Details', got %s", tool.Details)
		}
		if !slices.Equal(tool.UseCases, []string{"Test Use Case"}) {
			t.Errorf("Expected use cases to be ['Test Use Case'], got %v", tool.UseCases)
		}
		if !slices.Equal(tool.Tags, []string{"Test Tag"}) {
			t.Errorf("Expected tags to be ['Test Tag'], got %v", tool.Tags)
		}
		if tool.Website != "https://test.com" {
			t.Errorf("Expected website to be 'https://test.com', got %s", tool.Website)
		}
		if tool.Github != "https://github.com/test" {
			t.Errorf("Expected github to be 'https://github.com/test', got %s", tool.Github)
		}
		if tool.ReleaseYear != 2020 {
			t.Errorf("Expected release year to be 2020, got %d", tool.ReleaseYear)
		}
		if tool.Prolang != "Go" {
			t.Errorf("Expected pro lang to be 'Go', got %s", tool.Prolang)
		}
		if tool.Slug != "golang" {
			t.Errorf("Expected slug to be 'golang', got %s", tool.Slug)
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
			Devstatus: "invalid",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidToolDevstatus); !ok {
			t.Errorf("Expected ErrInvalidToolDevstatus, got %v", err)
		}
	})

	t.Run("Update tool with invalid fields", func(t *testing.T) {
		tool := testutil.CreateTestTool()

		err := tool.Update(domain.UpdateToolInput{
			Name:        "",
			Category:    "language",
			SubType:     "backend",
			Devstatus:   "active",
			Details:     "",
			UseCases:    []string{},
			Tags:        []string{},
			Website:     "not-a-valid-url",
			Github:      "not-a-valid-url",
			ReleaseYear: 1929,
			Prolang:     "",
			Slug:        "",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidField); ok {
			if _, exist := e.Fields["Name"]; !exist {
				t.Errorf("Expected the field \"Name\" to be present")
			}
			if _, exist := e.Fields["Slug"]; !exist {
				t.Errorf("Expected the field \"Slug\" to be present")
			}
			if _, exist := e.Fields["ReleaseYear"]; !exist {
				t.Errorf("Expected the field \"ReleaseYear\" to be present")
			}
			if _, exist := e.Fields["Prolang"]; !exist {
				t.Errorf("Expected the field \"Prolang\" to be present")
			}
			if _, exist := e.Fields["Website"]; !exist {
				t.Errorf("Expected the field \"Website\" to be present")
			}
			if _, exist := e.Fields["Github"]; !exist {
				t.Errorf("Expected the field \"Github\" to be present")
			}
		}
	})

	t.Run("Update tool with valid input", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		updateInput := domain.UpdateToolInput{
			Name:        "Mootools",
			Category:    "framework",
			SubType:     "frontend",
			Devstatus:   "deprecated",
			Details:     "Mootools is a framework",
			UseCases:    []string{"SPA", "SEO", "API", "Frontend"},
			Tags:        []string{"component-based", "declarative", "functional", "object-oriented"},
			Website:     "https://mootools.net",
			Github:      "https://github.com/mootools/mootools",
			ReleaseYear: 2006,
			Prolang:     "javascript",
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
		if tool.Devstatus != "deprecated" {
			t.Errorf("Expected dev status to be 'deprecated', got %s", tool.Devstatus)
		}
		if tool.Details != "Mootools is a framework" {
			t.Errorf("Expected details to be 'Mootools is a framework', got %s", tool.Details)
		}
		if !slices.Equal(tool.UseCases, []string{"SPA", "SEO", "API", "Frontend"}) {
			t.Errorf("Expected use cases to be ['SPA', 'SEO', 'API', 'Frontend'], got %v", tool.UseCases)
		}
		if !slices.Equal(tool.Tags, []string{"component-based", "declarative", "functional", "object-oriented"}) {
			t.Errorf("Expected tags to be ['component-based', 'declarative', 'functional', 'object-oriented'], got %v", tool.Tags)
		}
		if tool.Website != "https://mootools.net" {
			t.Errorf("Expected website to be 'https://mootools.net', got %s", tool.Website)
		}
		if tool.Github != "https://github.com/mootools/mootools" {
			t.Errorf("Expected github to be 'https://github.com/mootools/mootools', got %s", tool.Github)
		}
		if tool.ReleaseYear != 2006 {
			t.Errorf("Expected release year to be 2006, got %d", tool.ReleaseYear)
		}
		if tool.Prolang != "javascript" {
			t.Errorf("Expected pro lang to be 'javascript', got %s", tool.Prolang)
		}
		if tool.Slug != "mootools" {
			t.Errorf("Expected slug to be 'mootools', got %s", tool.Slug)
		}
	})
}
