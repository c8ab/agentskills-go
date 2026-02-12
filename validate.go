package agentskills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Validation limits per Agent Skills Spec.
const (
	MaxSkillNameLength     = 64
	MaxDescriptionLength   = 1024
	MaxCompatibilityLength = 500
)

// allowedFields defines the valid frontmatter fields per Agent Skills Spec.
var allowedFields = map[string]bool{
	"name":          true,
	"description":   true,
	"license":       true,
	"allowed-tools": true,
	"metadata":      true,
	"compatibility": true,
}

// validateName validates skill name format and directory match.
// Skill names support i18n characters (Unicode letters) plus hyphens.
// Names must be lowercase and cannot start/end with hyphens.
func validateName(name, skillDir string) []error {
	var errs []error

	if name == "" {
		errs = append(errs, ErrNameEmpty)
		return errs
	}

	// NFKC normalize the name
	name = norm.NFKC.String(strings.TrimSpace(name))

	if len(name) > MaxSkillNameLength {
		errs = append(errs, fmt.Errorf("skill name '%s' exceeds %d character limit (%d chars)",
			name, MaxSkillNameLength, len(name)))
	}

	if name != strings.ToLower(name) {
		errs = append(errs, fmt.Errorf("skill name '%s' must be lowercase", name))
	}

	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		errs = append(errs, ErrNameLeadingHyphen)
	}

	if strings.Contains(name, "--") {
		errs = append(errs, ErrNameConsecutiveHyphen)
	}

	// Check that all characters are alphanumeric or hyphen
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			errs = append(errs, fmt.Errorf("skill name '%s' contains invalid characters. Only letters, digits, and hyphens are allowed", name))
			break
		}
	}

	// Check directory name matches skill name
	if skillDir != "" {
		dirName := norm.NFKC.String(filepath.Base(skillDir))
		if dirName != name {
			errs = append(errs, fmt.Errorf("directory name '%s' must match skill name '%s'",
				filepath.Base(skillDir), name))
		}
	}

	return errs
}

// validateDescription validates description format.
func validateDescription(description string) []error {
	var errs []error

	if description == "" {
		errs = append(errs, ErrDescriptionEmpty)
		return errs
	}

	if len(description) > MaxDescriptionLength {
		errs = append(errs, fmt.Errorf("description exceeds %d character limit (%d chars)",
			MaxDescriptionLength, len(description)))
	}

	return errs
}

// validateCompatibility validates compatibility format.
func validateCompatibility(compatibility string) []error {
	var errs []error

	if len(compatibility) > MaxCompatibilityLength {
		errs = append(errs, fmt.Errorf("compatibility exceeds %d character limit (%d chars)",
			MaxCompatibilityLength, len(compatibility)))
	}

	return errs
}

// validateMetadataFields validates that only allowed fields are present.
func validateMetadataFields(metadata map[string]any) []error {
	var errs []error
	var extraFields []string

	for key := range metadata {
		if !allowedFields[key] {
			extraFields = append(extraFields, key)
		}
	}

	if len(extraFields) > 0 {
		allowed := make([]string, 0, len(allowedFields))
		for k := range allowedFields {
			allowed = append(allowed, k)
		}
		errs = append(errs, fmt.Errorf("unexpected fields in frontmatter: %s. Only %v are allowed",
			strings.Join(extraFields, ", "), allowed))
	}

	return errs
}

// ValidateMetadata validates parsed skill metadata.
// This is the core validation function that works on already-parsed metadata,
// avoiding duplicate file I/O when called from the parser.
func ValidateMetadata(metadata map[string]any, skillDir string) []error {
	var errs []error

	errs = append(errs, validateMetadataFields(metadata)...)

	name, hasName := metadata["name"]
	if !hasName {
		errs = append(errs, ErrMissingName)
	} else {
		nameStr, _ := name.(string)
		errs = append(errs, validateName(nameStr, skillDir)...)
	}

	desc, hasDesc := metadata["description"]
	if !hasDesc {
		errs = append(errs, ErrMissingDescription)
	} else {
		descStr, _ := desc.(string)
		errs = append(errs, validateDescription(descStr)...)
	}

	if compat, ok := metadata["compatibility"].(string); ok {
		errs = append(errs, validateCompatibility(compat)...)
	}

	return errs
}

// Validate validates a skill directory.
// Returns nil if valid, otherwise returns a ValidationErrors containing
// all problems found.
func Validate(skillDir string) error {
	info, err := os.Stat(skillDir)
	if os.IsNotExist(err) {
		return &ValidationErrors{Errors: []error{
			fmt.Errorf("path does not exist: %s", skillDir),
		}}
	}
	if err != nil {
		return &ValidationErrors{Errors: []error{err}}
	}

	if !info.IsDir() {
		return &ValidationErrors{Errors: []error{
			fmt.Errorf("not a directory: %s", skillDir),
		}}
	}

	skillMD := findSkillMD(skillDir)
	if skillMD == "" {
		return &ValidationErrors{Errors: []error{
			fmt.Errorf("missing required file: SKILL.md"),
		}}
	}

	content, err := os.ReadFile(skillMD)
	if err != nil {
		return &ValidationErrors{Errors: []error{err}}
	}

	metadata, _, err := parseFrontmatter(string(content))
	if err != nil {
		return &ValidationErrors{Errors: []error{err}}
	}

	errs := ValidateMetadata(metadata, skillDir)
	if len(errs) == 0 {
		return nil
	}

	return &ValidationErrors{Errors: errs}
}
