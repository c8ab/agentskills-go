// Package agentskills provides parsing and validation for Agent Skills.
package agentskills

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel errors for common validation conditions.
var (
	ErrSkillMDNotFound       = errors.New("SKILL.md not found")
	ErrMissingFrontmatter    = errors.New("SKILL.md must start with YAML frontmatter (---)")
	ErrUnclosedFrontmatter   = errors.New("SKILL.md frontmatter not properly closed with ---")
	ErrInvalidYAML           = errors.New("invalid YAML in frontmatter")
	ErrFrontmatterNotMapping = errors.New("SKILL.md frontmatter must be a YAML mapping")
	ErrMissingName           = errors.New("missing required field in frontmatter: name")
	ErrMissingDescription    = errors.New("missing required field in frontmatter: description")
	ErrNameEmpty             = errors.New("field 'name' must be a non-empty string")
	ErrDescriptionEmpty      = errors.New("field 'description' must be a non-empty string")
	ErrNameNotLowercase      = errors.New("skill name must be lowercase")
	ErrNameLeadingHyphen     = errors.New("skill name cannot start or end with a hyphen")
	ErrNameConsecutiveHyphen = errors.New("skill name cannot contain consecutive hyphens")
	ErrNameInvalidChars      = errors.New("skill name contains invalid characters")
	ErrPathNotExist          = errors.New("path does not exist")
	ErrPathNotDirectory      = errors.New("path is not a directory")
)

// ParseError indicates SKILL.md parsing failed.
type ParseError struct {
	Path string
	Err  error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s: %v", e.Path, e.Err)
	}
	return e.Err.Error()
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// ValidationError represents a single validation problem.
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// ValidationErrors holds multiple validation problems.
// It implements the error interface.
type ValidationErrors struct {
	Errors []error
}

func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	var sb strings.Builder
	sb.WriteString("validation failed with ")
	sb.WriteString(fmt.Sprintf("%d errors:\n", len(e.Errors)))
	for _, err := range e.Errors {
		sb.WriteString("  - ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Unwrap returns the list of errors for use with errors.Is/As.
func (e *ValidationErrors) Unwrap() []error {
	return e.Errors
}

// Add appends an error to the validation errors list.
func (e *ValidationErrors) Add(err error) {
	e.Errors = append(e.Errors, err)
}

// AddMessage appends a new error with the given message.
func (e *ValidationErrors) AddMessage(msg string) {
	e.Errors = append(e.Errors, errors.New(msg))
}

// HasErrors returns true if there are any validation errors.
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// AsError returns the ValidationErrors as an error, or nil if no errors.
func (e *ValidationErrors) AsError() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}
