# agentskills-go

[![CI](https://github.com/c8ab/agentskills-go/actions/workflows/ci.yml/badge.svg)](https://github.com/c8ab/agentskills-go/actions/workflows/ci.yml)

Go library for parsing and validating [Agent Skills](https://agentskills.dev) --
an open format for teaching AI agents specialized workflows via `SKILL.md` files.

## Installation

```bash
go get github.com/c8ab/agentskills-go
```

Requires Go 1.24 or later.

## Quick Start

### Read Skill Properties

```go
props, err := agentskills.ReadProperties("path/to/my-skill")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Skill: %s\n", props.Name)
fmt.Printf("Description: %s\n", props.Description)
```

### Validate a Skill

```go
if err := agentskills.Validate("path/to/my-skill"); err != nil {
    log.Fatal(err) // *ValidationErrors with all problems found
}
```

### Generate Agent Prompt

```go
prompt, err := agentskills.ToPrompt([]string{
    "path/to/skill-a",
    "path/to/skill-b",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(prompt)
```

Produces an `<available_skills>` XML block suitable for inclusion in LLM agent
system prompts:

```xml
<available_skills>
  <skill>
    <name>skill-a</name>
    <description>Description of skill A</description>
    <location>/absolute/path/to/skill-a/SKILL.md</location>
  </skill>
  <skill>
    <name>skill-b</name>
    <description>Description of skill B</description>
    <location>/absolute/path/to/skill-b/SKILL.md</location>
  </skill>
</available_skills>
```

## Documentation

See the [package documentation on pkg.go.dev](https://pkg.go.dev/github.com/c8ab/agentskills-go)
for the full API reference, including all exported types, functions, and error
types.

## License

[Apache-2.0](LICENSE)
