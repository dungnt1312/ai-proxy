# AI Proxy CLI

A unified command-line interface for multiple AI assistants (Claude, Kiro, Gemini, Cursor) with multi-agent workflow support.

## Features

- **Multi-Backend Support**: Seamlessly switch between Claude, Kiro, Gemini, and Cursor CLI
- **Multi-Agent Workflows**: Chain multiple AI agents for complex tasks (planning → review → coding → testing)
- **Project-Local Config**: Customize workflows per project
- **Context Awareness**: Automatically scans project structure before planning
- **Auto-Verify**: Runs build/test/lint after code changes
- **Diff-Based Review**: Reviews only changed files for faster feedback
- **Checkpoint/Resume**: Save progress and resume interrupted workflows

## Installation

```bash
# Clone and build
git clone <repo>
cd cli-proxy
go build -o proxy

# Install to PATH
cp proxy ~/.local/bin/ai-proxy
```

## Prerequisites

Install and authenticate at least one AI CLI:

```bash
# Claude Code CLI
npm install -g @anthropic-ai/claude-code

# Kiro CLI (AWS)
curl -fsSL https://cli.kiro.dev/install | bash

# Gemini CLI (Google)
npm install -g @google/gemini-cli

# Cursor Agent
curl -fsSL https://cursor.com/install | bash
```

Check installed CLIs:
```bash
claude --version       # Claude Code
kiro-cli --version     # Kiro CLI
gemini --version       # Gemini CLI
cursor-agent --version # Cursor Agent
```

## Quick Start

```bash
# Interactive mode
ai-proxy

# One-shot query
ai-proxy "explain golang interfaces"

# Use specific backend
ai-proxy -b kiro "what is docker?"

# Initialize project config
ai-proxy --init
```

## Commands

| Command | Description |
|---------|-------------|
| `/init` | Initialize project-local config (`.ai-proxy/config.json`) |
| `/switch <backend>` | Switch backend (claude, kiro, gemini, cursor) |
| `/list` | List available backends |
| `/workflow <name>` | Run a multi-agent workflow |
| `/resume [folder]` | Resume workflow (latest or specific folder) |
| `/skills` | List available skills |
| `/skill <name>` | Run a skill |
| `/skill install <url>` | Install skill from GitHub |
| `/skill remove <name>` | Remove a skill |
| `/skill info <name>` | Show skill details |
| `/clear` | Clear conversation history |
| `/help` | Show all commands |
| `quit` | Exit |

## CLI Flags

```bash
ai-proxy --init              # Initialize project config
ai-proxy -l                  # List backends
ai-proxy -b claude "hello"   # Use specific backend
ai-proxy --help              # Show help
```

## Workflows

### Built-in Workflows

| Workflow | Description | Stages |
|----------|-------------|--------|
| `feature` | Full feature development | plan → security → tasks → execute → verify → review → fix |
| `bugfix` | Bug analysis and fix | analyze → plan → fix |
| `refactor` | Code refactoring | analyze → review → execute |
| `api` | REST API development | plan → openapi → code → verify |
| `test` | Write tests | analyze → write → verify |
| `docs` | Generate documentation | scan → write |
| `docker` | Dockerize application | analyze → create → verify |

#### `feature` - Feature Development (7 stages)

```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ 1. PLAN  │ → │2.SECURITY│ → │ 3. TASKS │ → │4. EXECUTE│
│ (gemini) │   │  (kiro)  │   │  (kiro)  │   │ (claude) │
└──────────┘   └──────────┘   └──────────┘   └──────────┘
     ↓              ↓              ↓              ↓
 plan.md      security.md     tasks.md      [code files]

                                                 ↓

┌──────────┐   ┌──────────┐   ┌──────────┐
│ 7. FIX   │ ← │ 6.REVIEW │ ← │ 5.VERIFY │
│ (claude) │   │  (kiro)  │   │  (auto)  │
└──────────┘   └──────────┘   └──────────┘
     ↓              ↓              ↓
 [fix code]    review.md      verify.md
     │              ↑
     └──────────────┘  (loop if NEEDS_CHANGES)
```

#### `bugfix` - Bug Fix (3 stages)

```
analyze (gemini) → plan (kiro) → fix (claude)
```

#### `refactor` - Code Refactor (3 stages)

```
analyze (gemini) → review (kiro) → execute (claude)
```

#### `api` - REST API Development (4 stages)

```
plan (gemini) → openapi (kiro) → code (claude) → verify (auto)
```

#### `test` - Write Tests (3 stages)

```
analyze (gemini) → write (claude) → verify (auto)
```

#### `docs` - Generate Documentation (2 stages)

```
scan (gemini) → write (claude)
```

#### `docker` - Dockerize Application (3 stages)

```
analyze (gemini) → create (claude) → verify (kiro)
```

### Running Workflows

```bash
ai-proxy
[claude]> /workflow feature implement user authentication with JWT
[claude]> /workflow bugfix login returns 500 when user not found
[claude]> /workflow refactor split main.go into separate modules
[claude]> /workflow api create CRUD endpoints for products
[claude]> /workflow test write unit tests for auth module
[claude]> /workflow docs generate API documentation
[claude]> /workflow docker containerize the application
```

### Workflow Output

All workflow artifacts are saved to `.workflow/<timestamp>/`:

```
.workflow/
├── 20251216_230000/
│   ├── context.md      # Project context (auto-scanned)
│   ├── plan.md         # Implementation plan
│   ├── security.md     # Security review
│   ├── tasks.md        # Actionable tasks
│   ├── diff.md         # Changes made
│   ├── verify.md       # Build/test results
│   ├── review.md       # Code review
│   ├── state.json      # Checkpoint for resume
│   └── log.md          # Full workflow log
└── latest -> 20251216_230000/
```

## Configuration

### Global Config (`~/.ai-proxy.json`)

```json
{
  "default": "claude",
  "backends": {
    "claude": {
      "name": "Claude",
      "cmd": "claude",
      "args": [],
      "promptFlag": "-p",
      "resumeFlag": "--continue",
      "modelFlag": "--model"
    },
    "kiro": {
      "name": "Kiro",
      "cmd": "kiro-cli",
      "args": ["chat"],
      "promptFlag": "",
      "resumeFlag": "--resume",
      "modelFlag": "--model"
    },
    "gemini": {
      "name": "Gemini",
      "cmd": "gemini",
      "args": [],
      "promptFlag": "",
      "resumeFlag": "--resume",
      "modelFlag": "-m"
    },
    "cursor": {
      "name": "Cursor",
      "cmd": "cursor-agent",
      "args": [],
      "promptFlag": "-p",
      "resumeFlag": "--resume",
      "modelFlag": "--model"
    }
  }
}
```

### Project Config (`.ai-proxy/config.json`)

Initialize with `ai-proxy --init`, then customize:

```json
{
  "workflows": {
    "my-workflow": {
      "name": "My Custom Workflow",
      "stages": [
        {
          "name": "plan",
          "backend": "gemini",
          "model": "gemini-2.0-flash",
          "prompt": "Create a plan for: {{.Requirement}}",
          "outputFile": "plan.md",
          "interactive": false
        },
        {
          "name": "execute",
          "backend": "cursor",
          "model": "sonnet-4.5",
          "prompt": "Execute: {{.PlanContent}}",
          "outputFile": "",
          "interactive": true
        }
      ]
    }
  }
}
```

### Stage Configuration

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Stage identifier |
| `backend` | string | AI backend to use (claude, kiro, gemini, cursor) |
| `model` | string | Specific model (e.g., "opus", "sonnet-4.5", "gemini-2.0-flash") |
| `prompt` | string | Prompt template with variables |
| `outputFile` | string | Save output to this file (empty = no save) |
| `interactive` | bool | Run in interactive mode (for coding tasks) |
| `reviewLoop` | bool | Loop back if review fails |
| `maxAttempts` | int | Max review loop attempts before asking (default: 3) |

### Prompt Variables

| Variable | Description |
|----------|-------------|
| `{{.Requirement}}` | User's original requirement |
| `{{.ProjectContext}}` | Auto-scanned project info |
| `{{.PlanContent}}` | Content of plan.md |
| `{{.TasksContent}}` | Content of tasks.md |
| `{{.SecurityContent}}` | Content of security.md |
| `{{.DiffContent}}` | Content of diff.md |
| `{{.VerifyContent}}` | Content of verify.md |
| `{{.ReviewContent}}` | Content of review.md |

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              AI PROXY CLI                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Claude    │  │    Kiro     │  │   Gemini    │  │   Cursor    │        │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘        │
│         └─────────────────┼─────────────────┼───────────────┘              │
│                           │                 │                               │
│  ┌────────────────────────┴─────────────────┴────────────────────┐         │
│  │                     WORKFLOW ENGINE                            │         │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ │         │
│  │  │ Context │ │  Diff   │ │ Verify  │ │Checkpoint│ │  Init   │ │         │
│  │  │  Scan   │ │ Detect  │ │  Auto   │ │ Resume  │ │ Config  │ │         │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ │         │
│  └───────────────────────────────────────────────────────────────┘         │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────┐         │
│  │                         CONFIG                                 │         │
│  │  ~/.ai-proxy.json (global)  |  .ai-proxy/config.json (local)  │         │
│  └───────────────────────────────────────────────────────────────┘         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## File Structure

```
cli-proxy/
├── main.go         # REPL loop, backend calls
├── cmd.go          # CLI flags (cobra)
├── config.go       # Global config management
├── init.go         # Project-local config
├── workflow.go     # Workflow engine + definitions
├── context.go      # Project context scanning
├── diff.go         # File change detection
├── verify.go       # Auto build/test/vet
├── checkpoint.go   # Save/resume workflow state
├── utils.go        # Utilities (strip ANSI, etc.)
├── go.mod
└── go.sum
```

## Examples

### Simple Chat

```bash
ai-proxy
[claude]> explain golang interfaces
[claude]> /switch gemini
[gemini]> compare with typescript interfaces
```

### Feature Development

```bash
ai-proxy
[claude]> /workflow feature Create a REST API for todo app with CRUD operations

# Workflow runs:
# 1. Gemini creates implementation plan
# 2. Kiro reviews for security issues
# 3. Kiro creates actionable tasks
# 4. Claude implements the code (interactive)
# 5. Auto-verify runs go build/test
# 6. Kiro reviews the code changes
# 7. If issues found, Claude fixes them
# 8. Loop until approved
```

### Custom Workflow

```bash
# Create project config
ai-proxy --init

# Edit .ai-proxy/config.json to add custom workflow
# Then run it
ai-proxy
[claude]> /workflow my-custom-workflow implement feature X
```

### Resume Interrupted Workflow

```bash
# If workflow was interrupted (Ctrl+C, network error, etc.)
ai-proxy
[claude]> /resume
# Continues from last completed stage

# Resume specific workflow by folder name
[claude]> /resume 20251216_231305

# Or use 'latest' symlink
[claude]> /resume latest
```

## Skills

Skills are reusable prompt templates that can be run standalone or used in workflows.

### Built-in Skills

| Skill | Description | Backend |
|-------|-------------|---------|
| `code-review` | Review code changes for quality | kiro |
| `security-audit` | Audit code for security vulnerabilities | kiro |
| `explain` | Explain code or concepts in simple terms | gemini |
| `frontend-design` | Create distinctive, production-grade UI | claude |
| `mcp-builder` | Guide for creating MCP servers | claude |
| `skill-creator` | Create new skills | kiro |
| `project-skill-creator` | Analyze project and create project-specific skills | kiro |
| `git-commit` | Generate conventional commit messages | gemini |
| `refactor` | Code quality analysis and improvements | kiro |
| `test-generator` | Generate comprehensive unit tests | claude |

### Skill Structure

```
~/.ai-proxy/skills/
├── code-review/
│   ├── skill.yaml    # Metadata and config
│   └── prompt.md     # Prompt template
├── security-audit/
└── explain/
```

### skill.yaml Example

```yaml
name: code-review
description: Review code changes for quality
version: 1.0.0
author: ai-proxy

stage:
  backend: kiro
  model: ""
  interactive: false
  outputFile: review.md

inputs:
  - name: diff
    description: Code diff to review
    required: true
  - name: focus
    description: Areas to focus on
    required: false

tags: [review, quality]
```

### Using Skills

```bash
# List skills
/skills

# Run a skill
/skill code-review --diff="$(git diff)"
/skill frontend-design --requirement="Landing page for SaaS" --style=minimal
/skill git-commit --diff="$(git diff)"
/skill project-skill-creator --purpose="Create API endpoints following project patterns"

# Show skill info
/skill info code-review

# Install from GitHub
/skill install https://github.com/user/repo/tree/main/skills/my-skill

# Remove skill
/skill remove my-skill
```

### Skills in Workflows

Reference skills in workflow stages:

```json
{
  "workflows": {
    "my-workflow": {
      "stages": [
        {
          "name": "review",
          "skill": "code-review",
          "inputs": {
            "diff": "{{.DiffContent}}"
          }
        }
      ]
    }
  }
}
```

## License

MIT
