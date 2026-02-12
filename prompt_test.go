package agentskills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToPrompt_EmptyList(t *testing.T) {
	result, err := ToPrompt([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "<available_skills>\n</available_skills>"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestToPrompt_SingleSkill(t *testing.T) {
	result, err := ToPrompt([]string{"testdata/valid-skill"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "<available_skills>") {
		t.Error("expected <available_skills> tag")
	}
	if !strings.Contains(result, "</available_skills>") {
		t.Error("expected </available_skills> tag")
	}
	if !strings.Contains(result, "<name>\nvalid-skill\n</name>") {
		t.Error("expected skill name")
	}
	if !strings.Contains(result, "<description>\nA valid test skill\n</description>") {
		t.Error("expected skill description")
	}
	if !strings.Contains(result, "<location>") {
		t.Error("expected location tag")
	}
	if !strings.Contains(result, "SKILL.md") {
		t.Error("expected SKILL.md in location")
	}
}

func TestToPrompt_MultipleSkills(t *testing.T) {
	result, err := ToPrompt([]string{
		"testdata/valid-skill",
		"testdata/valid-all-fields",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Count skill tags
	skillCount := strings.Count(result, "<skill>")
	if skillCount != 2 {
		t.Errorf("expected 2 <skill> tags, got %d", skillCount)
	}

	closeSkillCount := strings.Count(result, "</skill>")
	if closeSkillCount != 2 {
		t.Errorf("expected 2 </skill> tags, got %d", closeSkillCount)
	}

	if !strings.Contains(result, "valid-skill") {
		t.Error("expected valid-skill in output")
	}
	if !strings.Contains(result, "valid-all-fields") {
		t.Error("expected valid-all-fields in output")
	}
}

func TestToPrompt_SpecialCharactersEscaped(t *testing.T) {
	result, err := ToPrompt([]string{"testdata/special-chars"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that special characters are escaped
	if !strings.Contains(result, "&lt;foo&gt;") {
		t.Error("expected <foo> to be escaped as &lt;foo&gt;")
	}
	if !strings.Contains(result, "&amp;") {
		t.Error("expected & to be escaped as &amp;")
	}
	if !strings.Contains(result, "&lt;bar&gt;") {
		t.Error("expected <bar> to be escaped as &lt;bar&gt;")
	}

	// Make sure raw special chars are not present
	if strings.Contains(result, "<foo>") {
		t.Error("raw <foo> should not be present")
	}
	if strings.Contains(result, "<bar>") {
		t.Error("raw <bar> should not be present")
	}
}

func TestToPrompt_InvalidSkillReturnsError(t *testing.T) {
	_, err := ToPrompt([]string{"testdata/nonexistent"})
	if err == nil {
		t.Error("expected error for nonexistent skill")
	}
}

func TestToPrompt_UsesAbsolutePath(t *testing.T) {
	result, err := ToPrompt([]string{"testdata/valid-skill"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The location should contain an absolute path
	cwd, _ := os.Getwd()
	expectedPath := filepath.Join(cwd, "testdata", "valid-skill", "SKILL.md")
	if !strings.Contains(result, expectedPath) {
		t.Errorf("expected absolute path %s in output, got: %s", expectedPath, result)
	}
}
