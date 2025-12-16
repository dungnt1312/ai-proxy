Create a new skill for ai-proxy.

## Purpose
{{.purpose}}

{{if .name}}
## Skill Name
{{.name}}
{{end}}

## Skill Structure

```
skill-name/
├── skill.yaml    # Metadata and config
└── prompt.md     # Prompt template
```

## skill.yaml Template

```yaml
name: skill-name
description: Clear description of what the skill does and when to use it
version: 1.0.0
author: your-name

stage:
  backend: kiro          # or claude, gemini, cursor
  model: ""              # optional specific model
  interactive: false     # true for coding tasks
  outputFile: output.md  # optional output file

inputs:
  - name: input_name
    description: What this input is for
    required: true
  - name: optional_input
    description: Optional parameter
    required: false
    default: default_value

tags: [relevant, tags, here]
```

## prompt.md Template

```markdown
Clear instruction for what to do.

## Input
{{.input_name}}

{{if .optional_input}}
## Optional Context
{{.optional_input}}
{{end}}

## Guidelines
- Specific instructions
- Best practices
- Output format

Provide the expected output.
```

## Core Principles

1. **Concise**: Only add context AI doesn't already have
2. **Clear Inputs**: Define required vs optional inputs
3. **Specific Output**: Define expected output format
4. **Appropriate Freedom**: Match specificity to task fragility

Generate the complete skill.yaml and prompt.md files.
