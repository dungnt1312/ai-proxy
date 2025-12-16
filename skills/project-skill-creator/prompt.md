Analyze the project and create a project-specific skill.

## Skill Purpose
{{.purpose}}

{{if .name}}
## Skill Name
{{.name}}
{{end}}

{{if .files}}
## Files to Analyze
{{.files}}
{{end}}

## Project Context
{{.ProjectContext}}

## Process

### Step 1: Analyze Project

1. **Read key files** to understand:
   - Project structure and architecture
   - Coding conventions and patterns
   - Naming conventions
   - Error handling patterns
   - Testing patterns
   - Common utilities/helpers

2. **Identify rules**:
   - File organization rules
   - Import/export patterns
   - Type definitions style
   - Comment/documentation style
   - Configuration patterns

### Step 2: Extract Patterns

Document the patterns found:
- How are similar components/modules structured?
- What boilerplate is repeated?
- What conventions must be followed?

### Step 3: Create Skill

Generate two files:

**skill.yaml**:
```yaml
name: project-specific-name
description: Clear description mentioning project-specific context
version: 1.0.0
author: project

stage:
  backend: claude
  model: ""
  interactive: true
  outputFile: ""

inputs:
  - name: relevant_input
    description: Based on what the skill needs
    required: true

tags: [project-specific, relevant-tags]
```

**prompt.md**:
```markdown
[Instructions that incorporate project-specific patterns]

## Project Conventions
- Convention 1 from analysis
- Convention 2 from analysis

## Template/Pattern
[Based on actual code patterns found]

## Guidelines
- Project-specific rules
- Naming conventions
- File locations
```

## Output

Provide:
1. Summary of patterns found
2. Complete skill.yaml content
3. Complete prompt.md content
4. Installation instructions:
   ```bash
   mkdir -p .ai-proxy/skills/skill-name
   # Save files to .ai-proxy/skills/skill-name/
   ```
