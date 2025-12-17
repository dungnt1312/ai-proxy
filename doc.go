// Package main implements AI Proxy CLI, a unified command-line interface that proxies prompts
// to multiple AI assistant CLIs (Claude, Kiro, Gemini, Cursor) and can orchestrate multi-stage
// workflows (plan → execute → verify → review → fix).
//
// This repository is primarily a CLI application (package main), but several exported types
// form the “public API” for workflow orchestration and skill execution:
//
//   - Workflow / Stage: Declarative workflow definitions with prompt templates and outputs.
//   - Skill: Reusable prompt templates loaded from skill folders (skill.yaml + prompt.md).
//   - FileSnapshot: Lightweight “diff” engine for listing changed/new files between snapshots.
//   - StageTimer: Simple stage timing + ETA estimation used while running workflows.
//
// For end-user usage (commands, flags, configuration file formats), see README.md.
package main
