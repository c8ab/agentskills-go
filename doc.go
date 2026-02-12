// Package agentskills provides parsing and validation for Agent Skills.
//
// Agent Skills is an open format for teaching AI agents specialized workflows.
// This package implements the core functionality for reading SKILL.md files,
// validating skill directories, and generating agent prompts.
//
// # Reading Skill Properties
//
// Use [ReadProperties] to parse the YAML frontmatter from a skill's SKILL.md file:
//
//	props, err := agentskills.ReadProperties("path/to/my-skill")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Skill: %s - %s\n", props.Name, props.Description)
//
// # Validating Skills
//
// Use [Validate] to check a skill directory against the Agent Skills specification:
//
//	err := agentskills.Validate("path/to/my-skill")
//	if err != nil {
//	    // err contains all validation problems
//	    log.Fatal(err)
//	}
//
// The validator checks:
//   - SKILL.md exists in the directory
//   - Valid YAML frontmatter with required fields (name, description)
//   - Name format: lowercase, max 64 chars, alphanumeric and hyphens only
//   - Name matches directory name
//   - Description: max 1024 chars
//   - No unexpected frontmatter fields
//
// # Generating Agent Prompts
//
// Use [ToPrompt] to generate the <available_skills> XML block for agent prompts:
//
//	prompt, err := agentskills.ToPrompt([]string{"path/to/skill-a", "path/to/skill-b"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Include prompt in agent system message
//
// For more information about Agent Skills, see https://agentskills.dev
package agentskills
