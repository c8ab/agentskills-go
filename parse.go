package agentskills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// findSkillMD finds the SKILL.md file in a skill directory.
// It prefers SKILL.md (uppercase) but accepts skill.md (lowercase).
// Returns the full path to the file, or empty string if not found.
func findSkillMD(skillDir string) string {
	for _, name := range []string{"SKILL.md", "skill.md"} {
		path := filepath.Join(skillDir, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// parseFrontmatter parses YAML frontmatter from SKILL.md content.
// Returns the parsed metadata map and the markdown body.
func parseFrontmatter(content string) (metadata map[string]any, body string, err error) {
	if !strings.HasPrefix(content, "---") {
		return nil, "", &ParseError{Err: ErrMissingFrontmatter}
	}

	// Split on "---" delimiter
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, "", &ParseError{Err: ErrUnclosedFrontmatter}
	}

	frontmatterStr := parts[1]
	body = strings.TrimSpace(parts[2])

	if err = yaml.Unmarshal([]byte(frontmatterStr), &metadata); err != nil {
		return nil, "", &ParseError{Err: fmt.Errorf("%w: %v", ErrInvalidYAML, err)}
	}

	if metadata == nil {
		return nil, "", &ParseError{Err: ErrFrontmatterNotMapping}
	}

	// Convert metadata values to strings where appropriate
	if rawMeta, ok := metadata["metadata"]; ok {
		if metaMap, ok := rawMeta.(map[string]any); ok {
			strMap := make(map[string]string)
			for k, v := range metaMap {
				strMap[k] = fmt.Sprintf("%v", v)
			}
			metadata["metadata"] = strMap
		}
	}

	return metadata, body, nil
}

// ReadProperties reads skill properties from SKILL.md frontmatter.
// This function parses the frontmatter and returns properties.
// It validates that required fields (name, description) exist but does NOT
// perform full validation. Use Validate() for complete validation.
func ReadProperties(skillDir string) (*SkillProperties, error) {
	skillMD := findSkillMD(skillDir)
	if skillMD == "" {
		return nil, &ParseError{Path: skillDir, Err: ErrSkillMDNotFound}
	}

	content, err := os.ReadFile(skillMD)
	if err != nil {
		return nil, &ParseError{Path: skillMD, Err: err}
	}

	metadata, _, err := parseFrontmatter(string(content))
	if err != nil {
		return nil, err
	}

	// Check required fields
	name, ok := metadata["name"]
	if !ok {
		return nil, &ValidationError{Field: "name", Err: ErrMissingName}
	}
	nameStr, ok := name.(string)
	if !ok || strings.TrimSpace(nameStr) == "" {
		return nil, &ValidationError{Field: "name", Err: ErrNameEmpty}
	}

	desc, ok := metadata["description"]
	if !ok {
		return nil, &ValidationError{Field: "description", Err: ErrMissingDescription}
	}
	descStr, ok := desc.(string)
	if !ok || strings.TrimSpace(descStr) == "" {
		return nil, &ValidationError{Field: "description", Err: ErrDescriptionEmpty}
	}

	props := &SkillProperties{
		Name:        strings.TrimSpace(nameStr),
		Description: strings.TrimSpace(descStr),
	}

	// Optional fields
	if v, ok := metadata["license"].(string); ok {
		props.License = v
	}
	if v, ok := metadata["compatibility"].(string); ok {
		props.Compatibility = v
	}
	if v, ok := metadata["allowed-tools"].(string); ok {
		props.AllowedTools = v
	}
	if v, ok := metadata["metadata"].(map[string]string); ok {
		props.Metadata = v
	}

	return props, nil
}
