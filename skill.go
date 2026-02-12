package agentskills

// SkillProperties holds parsed skill metadata from SKILL.md frontmatter.
type SkillProperties struct {
	// Name is the skill name in kebab-case (required).
	Name string `yaml:"name"`

	// Description explains what the skill does and when to use it (required).
	Description string `yaml:"description"`

	// License specifies the skill's license (optional).
	License string `yaml:"license,omitempty"`

	// Compatibility indicates environment requirements (optional).
	Compatibility string `yaml:"compatibility,omitempty"`

	// AllowedTools lists pre-approved tools the skill may use (optional, experimental).
	AllowedTools string `yaml:"allowed-tools,omitempty"`

	// Metadata holds arbitrary key-value pairs for client-specific properties (optional).
	Metadata map[string]string `yaml:"metadata,omitempty"`
}

// ToMap converts SkillProperties to a map, excluding empty values.
// The map uses the YAML field names (e.g., "allowed-tools" not "AllowedTools").
func (p *SkillProperties) ToMap() map[string]any {
	result := map[string]any{
		"name":        p.Name,
		"description": p.Description,
	}
	if p.License != "" {
		result["license"] = p.License
	}
	if p.Compatibility != "" {
		result["compatibility"] = p.Compatibility
	}
	if p.AllowedTools != "" {
		result["allowed-tools"] = p.AllowedTools
	}
	if len(p.Metadata) > 0 {
		result["metadata"] = p.Metadata
	}
	return result
}
