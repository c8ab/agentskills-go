package agentskills

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatter_Valid(t *testing.T) {
	content := `---
name: my-skill
description: A test skill
---
# My Skill

Instructions here.
`
	metadata, body, err := parseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata["name"] != "my-skill" {
		t.Errorf("expected name 'my-skill', got %v", metadata["name"])
	}
	if metadata["description"] != "A test skill" {
		t.Errorf("expected description 'A test skill', got %v", metadata["description"])
	}
	if body != "# My Skill\n\nInstructions here." {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestParseFrontmatter_MissingFrontmatter(t *testing.T) {
	content := "# No frontmatter here"
	_, _, err := parseFrontmatter(content)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrMissingFrontmatter) {
		t.Errorf("expected ErrMissingFrontmatter, got %v", err)
	}
}

func TestParseFrontmatter_UnclosedFrontmatter(t *testing.T) {
	content := `---
name: my-skill
description: A test skill
`
	_, _, err := parseFrontmatter(content)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrUnclosedFrontmatter) {
		t.Errorf("expected ErrUnclosedFrontmatter, got %v", err)
	}
}

func TestParseFrontmatter_InvalidYAML(t *testing.T) {
	content := `---
name: [invalid
description: broken
---
Body here
`
	_, _, err := parseFrontmatter(content)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidYAML) {
		t.Errorf("expected ErrInvalidYAML, got %v", err)
	}
}

func TestReadProperties_ValidSkill(t *testing.T) {
	props, err := ReadProperties("testdata/valid-skill")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if props.Name != "valid-skill" {
		t.Errorf("expected name 'valid-skill', got %q", props.Name)
	}
	if props.Description != "A valid test skill" {
		t.Errorf("expected description 'A valid test skill', got %q", props.Description)
	}
}

func TestReadProperties_AllFields(t *testing.T) {
	props, err := ReadProperties("testdata/valid-all-fields")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if props.Name != "valid-all-fields" {
		t.Errorf("expected name 'valid-all-fields', got %q", props.Name)
	}
	if props.License != "MIT" {
		t.Errorf("expected license 'MIT', got %q", props.License)
	}
	if props.Compatibility != "Requires Python 3.11+" {
		t.Errorf("expected compatibility 'Requires Python 3.11+', got %q", props.Compatibility)
	}
	if props.AllowedTools != "Bash(git:*) Bash(jq:*)" {
		t.Errorf("expected allowed-tools 'Bash(git:*) Bash(jq:*)', got %q", props.AllowedTools)
	}
	if props.Metadata["author"] != "Test Author" {
		t.Errorf("expected metadata.author 'Test Author', got %q", props.Metadata["author"])
	}
	if props.Metadata["version"] != "1.0" {
		t.Errorf("expected metadata.version '1.0', got %q", props.Metadata["version"])
	}
}

func TestReadProperties_MissingSkillMD(t *testing.T) {
	_, err := ReadProperties("testdata/nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrSkillMDNotFound) {
		t.Errorf("expected ErrSkillMDNotFound, got %v", err)
	}
}

func TestReadProperties_MissingName(t *testing.T) {
	_, err := ReadProperties("testdata/missing-name")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("expected ErrMissingName, got %v", err)
	}
}

func TestReadProperties_MissingDescription(t *testing.T) {
	_, err := ReadProperties("testdata/missing-description")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrMissingDescription) {
		t.Errorf("expected ErrMissingDescription, got %v", err)
	}
}

func TestFindSkillMD_PrefersUppercase(t *testing.T) {
	// Create temp dir with both SKILL.md and skill.md
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "SKILL.md"), []byte("uppercase"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "skill.md"), []byte("lowercase"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := findSkillMD(tmpDir)
	if filepath.Base(result) != "SKILL.md" {
		t.Errorf("expected SKILL.md, got %s", filepath.Base(result))
	}
}

func TestFindSkillMD_AcceptsLowercase(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "skill.md"), []byte("lowercase"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := findSkillMD(tmpDir)
	if result == "" {
		t.Error("expected to find skill.md, got empty string")
	}
}

func TestFindSkillMD_ReturnsEmptyWhenMissing(t *testing.T) {
	tmpDir := t.TempDir()
	result := findSkillMD(tmpDir)
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestSkillProperties_ToMap(t *testing.T) {
	props := &SkillProperties{
		Name:         "test-skill",
		Description:  "A test",
		License:      "MIT",
		AllowedTools: "Bash(git:*)",
		Metadata:     map[string]string{"author": "test"},
	}

	m := props.ToMap()

	if m["name"] != "test-skill" {
		t.Errorf("expected name 'test-skill', got %v", m["name"])
	}
	if m["description"] != "A test" {
		t.Errorf("expected description 'A test', got %v", m["description"])
	}
	if m["license"] != "MIT" {
		t.Errorf("expected license 'MIT', got %v", m["license"])
	}
	if m["allowed-tools"] != "Bash(git:*)" {
		t.Errorf("expected allowed-tools 'Bash(git:*)', got %v", m["allowed-tools"])
	}
	if _, ok := m["compatibility"]; ok {
		t.Error("expected compatibility to be omitted when empty")
	}
}
