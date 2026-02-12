package agentskills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidate_ValidSkill(t *testing.T) {
	err := Validate("testdata/valid-skill")
	if err != nil {
		t.Errorf("expected no error for valid skill, got: %v", err)
	}
}

func TestValidate_ValidAllFields(t *testing.T) {
	err := Validate("testdata/valid-all-fields")
	if err != nil {
		t.Errorf("expected no error for valid skill with all fields, got: %v", err)
	}
}

func TestValidate_NonexistentPath(t *testing.T) {
	err := Validate("testdata/nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("expected 'does not exist' error, got: %v", err)
	}
}

func TestValidate_NotADirectory(t *testing.T) {
	err := Validate("testdata/valid-skill/SKILL.md")
	if err == nil {
		t.Fatal("expected error for file path")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("expected 'not a directory' error, got: %v", err)
	}
}

func TestValidate_MissingSkillMD(t *testing.T) {
	tmpDir := t.TempDir()
	err := Validate(tmpDir)
	if err == nil {
		t.Fatal("expected error for missing SKILL.md")
	}
	if !strings.Contains(err.Error(), "SKILL.md") {
		t.Errorf("expected 'SKILL.md' in error, got: %v", err)
	}
}

func TestValidate_InvalidNameUppercase(t *testing.T) {
	err := Validate("testdata/invalid-uppercase")
	if err == nil {
		t.Fatal("expected error for uppercase name")
	}
	if !strings.Contains(err.Error(), "lowercase") {
		t.Errorf("expected 'lowercase' error, got: %v", err)
	}
}

func TestValidate_NameTooLong(t *testing.T) {
	tmpDir := t.TempDir()
	longName := strings.Repeat("a", 70)
	skillDir := filepath.Join(tmpDir, longName)
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "---\nname: " + longName + "\ndescription: A test skill\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err == nil {
		t.Fatal("expected error for name too long")
	}
	if !strings.Contains(err.Error(), "exceeds") && !strings.Contains(err.Error(), "character limit") {
		t.Errorf("expected 'exceeds character limit' error, got: %v", err)
	}
}

func TestValidate_NameLeadingHyphen(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "-my-skill")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "---\nname: -my-skill\ndescription: A test skill\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err == nil {
		t.Fatal("expected error for leading hyphen")
	}
	if !strings.Contains(err.Error(), "start or end with a hyphen") {
		t.Errorf("expected hyphen error, got: %v", err)
	}
}

func TestValidate_NameConsecutiveHyphens(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "my--skill")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "---\nname: my--skill\ndescription: A test skill\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err == nil {
		t.Fatal("expected error for consecutive hyphens")
	}
	if !strings.Contains(err.Error(), "consecutive hyphens") {
		t.Errorf("expected consecutive hyphens error, got: %v", err)
	}
}

func TestValidate_NameInvalidCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "my_skill")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "---\nname: my_skill\ndescription: A test skill\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err == nil {
		t.Fatal("expected error for invalid characters")
	}
	if !strings.Contains(err.Error(), "invalid characters") {
		t.Errorf("expected invalid characters error, got: %v", err)
	}
}

func TestValidate_NameDirectoryMismatch(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "wrong-name")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "---\nname: correct-name\ndescription: A test skill\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err == nil {
		t.Fatal("expected error for name/directory mismatch")
	}
	if !strings.Contains(err.Error(), "must match skill name") {
		t.Errorf("expected 'must match skill name' error, got: %v", err)
	}
}

func TestValidate_UnexpectedFields(t *testing.T) {
	err := Validate("testdata/unexpected-fields")
	if err == nil {
		t.Fatal("expected error for unexpected fields")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "unexpected fields") {
		t.Errorf("expected 'unexpected fields' error, got: %v", err)
	}
}

func TestValidate_DescriptionTooLong(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "my-skill")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	longDesc := strings.Repeat("x", 1100)
	content := "---\nname: my-skill\ndescription: " + longDesc + "\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err == nil {
		t.Fatal("expected error for description too long")
	}
	if !strings.Contains(err.Error(), "1024") {
		t.Errorf("expected '1024' in error, got: %v", err)
	}
}

func TestValidate_CompatibilityTooLong(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "my-skill")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	longCompat := strings.Repeat("x", 550)
	content := "---\nname: my-skill\ndescription: A test skill\ncompatibility: " + longCompat + "\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err == nil {
		t.Fatal("expected error for compatibility too long")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected '500' in error, got: %v", err)
	}
}

func TestValidate_I18NChineseName(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "技能")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "---\nname: 技能\ndescription: A skill with Chinese name\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err != nil {
		t.Errorf("expected no error for Chinese name, got: %v", err)
	}
}

func TestValidate_I18NRussianLowercase(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "навык")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "---\nname: навык\ndescription: A skill with Russian lowercase name\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err != nil {
		t.Errorf("expected no error for Russian lowercase name, got: %v", err)
	}
}

func TestValidate_I18NRussianUppercaseRejected(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "НАВЫК")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "---\nname: НАВЫК\ndescription: A skill with Russian uppercase name\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err == nil {
		t.Fatal("expected error for Russian uppercase name")
	}
	if !strings.Contains(err.Error(), "lowercase") {
		t.Errorf("expected 'lowercase' error, got: %v", err)
	}
}

func TestValidate_NFKCNormalization(t *testing.T) {
	tmpDir := t.TempDir()
	// Use composed form for directory name
	composedName := "café"
	skillDir := filepath.Join(tmpDir, composedName)
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Use decomposed form in SKILL.md (e + combining acute)
	decomposedName := "cafe\u0301"
	content := "---\nname: " + decomposedName + "\ndescription: A test skill\n---\nBody\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Validate(skillDir)
	if err != nil {
		t.Errorf("expected no error after NFKC normalization, got: %v", err)
	}
}

func TestValidationErrors_Interface(t *testing.T) {
	errs := &ValidationErrors{}
	errs.AddMessage("error 1")
	errs.AddMessage("error 2")

	if !errs.HasErrors() {
		t.Error("expected HasErrors to return true")
	}

	errStr := errs.Error()
	if !strings.Contains(errStr, "error 1") || !strings.Contains(errStr, "error 2") {
		t.Errorf("expected both errors in string, got: %s", errStr)
	}

	if errs.AsError() == nil {
		t.Error("expected AsError to return non-nil")
	}

	emptyErrs := &ValidationErrors{}
	if emptyErrs.HasErrors() {
		t.Error("expected HasErrors to return false for empty")
	}
	if emptyErrs.AsError() != nil {
		t.Error("expected AsError to return nil for empty")
	}
}
