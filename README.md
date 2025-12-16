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
# Claude CLI
npm install -g @anthropic-ai/claude-code
claude login

# Kiro CLI (AWS)
# Follow AWS Kiro installation guide

# Gemini CLI
npm install -g @anthropic-ai/gemini-cli
gemini login

# Cursor Agent
npm install -g cursor-agent
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
| `/resume` | Resume last interrupted workflow |
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

### Running Workflows

```bash
ai-proxy
[claude]> /workflow feature implement user authentication with JWT
[claude]> /workflow bugfix login returns 500 when user not found
[claude]> /workflow refactor split main.go into separate modules
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
```

## License

MIT
