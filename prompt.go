package agentskills

import (
	"html"
	"path/filepath"
	"strings"
)

// ToPrompt generates the <available_skills> XML block for inclusion in agent prompts.
//
// This XML format is what Anthropic uses and recommends for Claude models.
// Skill Clients may format skill information differently to suit their
// models or preferences.
//
// Example output:
//
//	<available_skills>
//	<skill>
//	<name>pdf-reader</name>
//	<description>Read and extract text from PDF files</description>
//	<location>/path/to/pdf-reader/SKILL.md</location>
//	</skill>
//	</available_skills>
func ToPrompt(skillDirs []string) (string, error) {
	if len(skillDirs) == 0 {
		return "<available_skills>\n</available_skills>", nil
	}

	var lines []string
	lines = append(lines, "<available_skills>")

	for _, skillDir := range skillDirs {
		absDir, err := filepath.Abs(skillDir)
		if err != nil {
			return "", err
		}

		props, err := ReadProperties(absDir)
		if err != nil {
			return "", err
		}

		skillMDPath := findSkillMD(absDir)

		lines = append(lines, "<skill>")
		lines = append(lines, "<name>")
		lines = append(lines, html.EscapeString(props.Name))
		lines = append(lines, "</name>")
		lines = append(lines, "<description>")
		lines = append(lines, html.EscapeString(props.Description))
		lines = append(lines, "</description>")
		lines = append(lines, "<location>")
		lines = append(lines, skillMDPath)
		lines = append(lines, "</location>")
		lines = append(lines, "</skill>")
	}

	lines = append(lines, "</available_skills>")

	return strings.Join(lines, "\n"), nil
}
