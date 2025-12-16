package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Stage struct {
	Name        string `json:"name"`
	Backend     string `json:"backend"`
	Model       string `json:"model,omitempty"`
	Prompt      string `json:"prompt"`
	OutputFile  string `json:"outputFile"`
	Interactive bool   `json:"interactive"`
	ReviewLoop  bool   `json:"reviewLoop"`
	Skippable   bool   `json:"skippable,omitempty"`
}

type Workflow struct {
	Name   string  `json:"name"`
	Stages []Stage `json:"stages"`
}

type WorkflowContext struct {
	Requirement    string
	WorkDir        string
	Results        map[string]string
	CurrentIdx     int
	LogFile        *os.File
	BeforeSnapshot *FileSnapshot
}

var defaultWorkflows = map[string]Workflow{
	"feature": {
		Name: "Feature Development",
		Stages: []Stage{
			{
				Name:       "plan",
				Backend:    "gemini",
				OutputFile: "plan.md",
				Prompt: `You are a software architect. Create a detailed implementation plan.

{{.ProjectContext}}

## Requirement
{{.Requirement}}

Output in Markdown:
# Implementation Plan

## Overview
Brief description

## Files to Create/Modify
- file1.go - description

## Dependencies
- list any NEW packages needed (check existing in context)

## Implementation Steps
1. Step one
2. Step two

Be thorough but concise. Consider existing project structure.`,
			},
			{
				Name:       "security",
				Backend:    "kiro",
				OutputFile: "security.md",
				Skippable:  true,
				Prompt: `Review this plan for security issues:

{{.PlanContent}}

Requirement: {{.Requirement}}

Output:
# Security Review
## Risk: LOW/MEDIUM/HIGH
## Concerns
- Issue and mitigation
## Recommendations
- Best practices`,
			},
			{
				Name:       "tasks",
				Backend:    "kiro",
				OutputFile: "tasks.md",
				Prompt: `Create tasks from this plan:

{{.PlanContent}}

Security notes:
{{.SecurityContent}}

Requirement: {{.Requirement}}

Output:
# Tasks
## Task 1: [Title]
- **File:** filename
- **Action:** create/modify
- **Description:** what to do
- **Security:** notes if any`,
			},
			{
				Name:        "execute",
				Backend:     "claude",
				Interactive: true,
				Prompt: `Execute these tasks to implement the feature:

{{.TasksContent}}

Original requirement: {{.Requirement}}

Work through each task. Create/modify files as specified.
When done, type /quit or exit to continue to verification.`,
			},
			{
				Name:       "verify",
				Backend:    "auto",
				OutputFile: "verify.md",
				Prompt:     "", // Special stage, no LLM call
			},
			{
				Name:       "code-review",
				Backend:    "kiro",
				OutputFile: "review.md",
				Prompt: `You are a senior code reviewer. Review the code changes:

## Requirement
{{.Requirement}}

## Changes Made
{{.DiffContent}}

## Verification Results
{{.VerifyContent}}

Review the changes and output:
# Code Review

## Status: APPROVED / NEEDS_CHANGES

## Summary
Brief overview

## Issues Found (if any)
- Issue 1: file, line, description, fix

## Suggestions
- Optional improvements

If code is good and tests pass, say "Status: APPROVED".`,
			},
			{
				Name:        "fix",
				Backend:     "claude",
				Interactive: true,
				ReviewLoop:  true,
				Prompt: `Fix these code review issues:

{{.ReviewContent}}

Original requirement: {{.Requirement}}

Address each issue listed. When done, exit to continue review.`,
			},
		},
	},
	"bugfix": {
		Name: "Bug Fix",
		Stages: []Stage{
			{
				Name:       "analyze",
				Backend:    "gemini",
				OutputFile: "analysis.md",
				Prompt: `You are a debugging expert. Analyze this bug:
{{.Requirement}}

Output in Markdown:
# Bug Analysis

## Problem Summary
What's happening

## Possible Causes
1. Cause one
2. Cause two

## Files to Investigate
- file1.go - why

## Suggested Fix
Step by step solution`,
			},
			{
				Name:       "plan",
				Backend:    "kiro",
				OutputFile: "fix-tasks.md",
				Prompt: `Review this bug analysis and create fix tasks:

{{.PlanContent}}

Original bug: {{.Requirement}}

Output:
# Fix Tasks

## Task 1: [Title]
- **File:** filename
- **Change:** what to change
- **Code:** before/after snippet`,
			},
			{
				Name:        "fix",
				Backend:     "claude",
				Interactive: true,
				Prompt: `Fix this bug following these tasks:

{{.TasksContent}}

Original bug: {{.Requirement}}

Apply the fixes carefully.`,
			},
		},
	},
	"refactor": {
		Name: "Code Refactor",
		Stages: []Stage{
			{
				Name:       "analyze",
				Backend:    "gemini",
				OutputFile: "refactor-plan.md",
				Prompt: `You are a code quality expert. Plan a refactor for:
{{.Requirement}}

Output:
# Refactor Plan

## Current Issues
- issue 1

## Proposed Changes
1. Change one

## Files Affected
- file1.go

## Risk Assessment
Low/Medium/High and why`,
			},
			{
				Name:       "review",
				Backend:    "kiro",
				OutputFile: "refactor-tasks.md",
				Prompt: `Review this refactor plan and create safe tasks:

{{.PlanContent}}

Original request: {{.Requirement}}

Output:
# Refactor Tasks

## Task 1: [Title]
- **File:** filename
- **Change:** description
- **Before:** old code
- **After:** new code`,
			},
			{
				Name:        "execute",
				Backend:     "claude",
				Interactive: true,
				Prompt: `Execute this refactor:

{{.TasksContent}}

Original request: {{.Requirement}}

Make changes carefully.`,
			},
		},
	},
	"api": {
		Name: "REST API Development",
		Stages: []Stage{
			{
				Name:       "plan",
				Backend:    "gemini",
				OutputFile: "api-plan.md",
				Prompt: `Design a REST API for: {{.Requirement}}

{{.ProjectContext}}

Output:
# API Design
## Endpoints
- METHOD /path - description
## Data Models
## Authentication
## Error Handling`,
			},
			{
				Name:       "openapi",
				Backend:    "kiro",
				OutputFile: "openapi.yaml",
				Skippable:  true,
				Prompt: `Create OpenAPI 3.0 spec from this design:

{{.PlanContent}}

Output valid YAML only.`,
			},
			{
				Name:        "code",
				Backend:     "claude",
				Interactive: true,
				Prompt: `Implement this API:

{{.PlanContent}}

Requirement: {{.Requirement}}`,
			},
			{
				Name:       "verify",
				Backend:    "auto",
				OutputFile: "verify.md",
			},
		},
	},
	"test": {
		Name: "Write Tests",
		Stages: []Stage{
			{
				Name:       "analyze",
				Backend:    "gemini",
				OutputFile: "test-plan.md",
				Prompt: `Analyze code and plan tests for: {{.Requirement}}

{{.ProjectContext}}

Output:
# Test Plan
## Files to Test
## Test Cases
- test case 1
- test case 2
## Edge Cases`,
			},
			{
				Name:        "write",
				Backend:     "claude",
				Interactive: true,
				Prompt: `Write tests based on this plan:

{{.PlanContent}}

Requirement: {{.Requirement}}

Create comprehensive unit tests.`,
			},
			{
				Name:       "verify",
				Backend:    "auto",
				OutputFile: "verify.md",
			},
		},
	},
	"docs": {
		Name: "Generate Documentation",
		Stages: []Stage{
			{
				Name:       "scan",
				Backend:    "gemini",
				OutputFile: "doc-outline.md",
				Prompt: `Analyze project and create documentation outline:

{{.ProjectContext}}

Requirement: {{.Requirement}}

Output:
# Documentation Outline
## Overview
## Installation
## Usage
## API Reference
## Examples`,
			},
			{
				Name:        "write",
				Backend:     "claude",
				Interactive: true,
				Prompt: `Write documentation based on this outline:

{{.PlanContent}}

Requirement: {{.Requirement}}

Create clear, comprehensive docs.`,
			},
		},
	},
	"docker": {
		Name: "Dockerize Application",
		Stages: []Stage{
			{
				Name:       "analyze",
				Backend:    "gemini",
				OutputFile: "docker-plan.md",
				Prompt: `Analyze project for containerization:

{{.ProjectContext}}

Requirement: {{.Requirement}}

Output:
# Docker Plan
## Base Image
## Dependencies
## Build Steps
## Ports
## Environment Variables
## Volumes`,
			},
			{
				Name:        "create",
				Backend:     "claude",
				Interactive: true,
				Prompt: `Create Docker configuration:

{{.PlanContent}}

Requirement: {{.Requirement}}

Create Dockerfile and docker-compose.yml if needed.`,
			},
			{
				Name:       "verify",
				Backend:    "kiro",
				OutputFile: "docker-review.md",
				Prompt: `Review Docker configuration for:
- Security best practices
- Image size optimization
- Multi-stage builds
- Proper layer caching

{{.DiffContent}}

Output issues and suggestions.`,
			},
		},
	},
}

func (wf *Workflow) Run(requirement string) error {
	baseDir := ".workflow"
	timestamp := time.Now().Format("20060102_150405")
	workDir := filepath.Join(baseDir, timestamp)
	os.MkdirAll(workDir, 0755)

	latestLink := filepath.Join(baseDir, "latest")
	os.Remove(latestLink)
	os.Symlink(timestamp, latestLink)

	logPath := filepath.Join(workDir, "log.md")
	logFile, _ := os.Create(logPath)
	defer logFile.Close()

	ctx := &WorkflowContext{
		Requirement:    requirement,
		WorkDir:        workDir,
		Results:        make(map[string]string),
		LogFile:        logFile,
		BeforeSnapshot: takeSnapshot(),
	}

	ctx.log("# Workflow: %s\n", wf.Name)
	ctx.log("**Requirement:** %s\n", requirement)
	ctx.log("**Time:** %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	projectCtx := scanProjectContext()
	ctx.Results["project-context"] = projectCtx
	os.WriteFile(filepath.Join(workDir, "context.md"), []byte(projectCtx), 0644)

	fmt.Printf("\n%s Workflow: %s\n", cyan("▶"), wf.Name)
	fmt.Printf("%s Requirement: %s\n", dim("│"), requirement)
	fmt.Printf("%s Directory: %s\n", dim("│"), workDir)
	fmt.Printf("%s Context: scanned project\n\n", dim("│"))

	i := 0
	maxReviewLoops := 3
	reviewLoopCount := 0

	for i < len(wf.Stages) {
		stage := wf.Stages[i]
		ctx.CurrentIdx = i

		if stage.ReviewLoop && reviewLoopCount == 0 {
			reviewContent := ctx.Results["code-review"]
			if strings.Contains(strings.ToUpper(reviewContent), "APPROVED") &&
				!strings.Contains(strings.ToUpper(reviewContent), "NEEDS_CHANGES") {
				fmt.Printf("%s Skipping fix (code approved)\n\n", green("✓"))
				i++
				continue
			}
		}

		fmt.Printf("%s [Stage %d/%d] %s (%s)\n", cyan("●"), i+1, len(wf.Stages), stage.Name, stage.Backend)

		if stage.OutputFile != "" {
			fmt.Printf("%s Output: %s\n", dim("│"), filepath.Join(workDir, stage.OutputFile))
		}

		if stage.Skippable {
			fmt.Printf("%s Skip this stage? [y/N]: ", yellow("?"))
			var input string
			fmt.Scanln(&input)
			if strings.ToLower(strings.TrimSpace(input)) == "y" {
				fmt.Printf("%s Skipped\n\n", dim("○"))
				i++
				continue
			}
		}

		ctx.log("## Stage %d: %s (loop %d)\n\n", i+1, stage.Name, reviewLoopCount)

		var result string
		var err error

		if stage.Backend == "auto" && stage.Name == "verify" {
			afterSnapshot := takeSnapshot()
			diffContent := ctx.BeforeSnapshot.Diff(afterSnapshot)
			ctx.Results["diff"] = diffContent
			os.WriteFile(filepath.Join(workDir, "diff.md"), []byte(diffContent), 0644)
			fmt.Printf("%s Generated diff of changes\n", dim("│"))

			verifyOutput, passed := runVerifyStage(workDir)
			result = verifyOutput
			ctx.Results["verify"] = result
			if !passed {
				fmt.Printf("%s Some checks failed, will be included in review\n", yellow("!"))
			}
		} else {
			result, err = runStage(&stage, ctx)
			if err != nil {
				return fmt.Errorf("stage %s failed: %w", stage.Name, err)
			}
		}

		ctx.Results[stage.Name] = result

		if stage.OutputFile != "" {
			outPath := filepath.Join(workDir, stage.OutputFile)
			cleanResult := stripANSI(result)
			os.WriteFile(outPath, []byte(cleanResult), 0644)
			fmt.Printf("%s Saved: %s\n", green("✓"), outPath)
		}

		ctx.log("### Output\n```\n%s\n```\n\n", truncate(result, 2000))
		fmt.Printf("%s Stage completed\n\n", green("✓"))

		saveCheckpoint(ctx, wf.Name, i)

		if stage.Name == "code-review" {
			if strings.Contains(strings.ToUpper(result), "NEEDS_CHANGES") {
				reviewLoopCount++
				if reviewLoopCount >= maxReviewLoops {
					fmt.Printf("%s Max review loops reached (%d)\n", yellow("!"), maxReviewLoops)
				} else {
					fmt.Printf("%s Changes needed, going to fix stage (loop %d/%d)\n\n", yellow("↻"), reviewLoopCount, maxReviewLoops)
					i++ // Go to fix stage
					continue
				}
			} else {
				fmt.Printf("%s Code approved!\n\n", green("✓"))
			}
		}

		if stage.ReviewLoop && reviewLoopCount > 0 && reviewLoopCount < maxReviewLoops {
			for j, s := range wf.Stages {
				if s.Name == "code-review" {
					i = j
					fmt.Printf("%s Back to code review\n\n", cyan("↻"))
					continue
				}
			}
		}

		i++
	}

	fmt.Printf("%s Workflow completed!\n", green("✓"))
	fmt.Printf("%s Files in: %s/\n", dim("📁"), workDir)

	files, _ := os.ReadDir(workDir)
	for _, f := range files {
		fmt.Printf("   %s %s\n", dim("•"), f.Name())
	}

	return nil
}

func resumeWorkflow() error {
	workDir := findLatestWorkflow()
	if workDir == "" {
		return fmt.Errorf("no workflow to resume")
	}

	state, err := loadCheckpoint(workDir)
	if err != nil {
		return fmt.Errorf("cannot load checkpoint: %w", err)
	}

	wf := getWorkflow(state.WorkflowName)
	if wf == nil {
		return fmt.Errorf("unknown workflow: %s", state.WorkflowName)
	}

	fmt.Printf("%s Resuming: %s (stage %d/%d)\n", cyan("↻"), wf.Name, state.CurrentStage+2, len(wf.Stages))
	fmt.Printf("%s Directory: %s\n\n", dim("│"), workDir)

	logFile, _ := os.OpenFile(filepath.Join(workDir, "log.md"), os.O_APPEND|os.O_WRONLY, 0644)
	defer logFile.Close()

	ctx := &WorkflowContext{
		Requirement:    state.Requirement,
		WorkDir:        workDir,
		Results:        state.Results,
		LogFile:        logFile,
		BeforeSnapshot: takeSnapshot(),
	}

	for i := state.CurrentStage + 1; i < len(wf.Stages); i++ {
		stage := wf.Stages[i]
		fmt.Printf("%s [Stage %d/%d] %s (%s)\n", cyan("●"), i+1, len(wf.Stages), stage.Name, stage.Backend)

		result, _ := runStage(&stage, ctx)
		ctx.Results[stage.Name] = result

		if stage.OutputFile != "" {
			outPath := filepath.Join(workDir, stage.OutputFile)
			os.WriteFile(outPath, []byte(stripANSI(result)), 0644)
		}
		fmt.Printf("%s Stage completed\n\n", green("✓"))
		saveCheckpoint(ctx, wf.Name, i)
	}

	fmt.Printf("%s Workflow resumed and completed!\n", green("✓"))
	return nil
}

func (ctx *WorkflowContext) log(format string, args ...interface{}) {
	if ctx.LogFile != nil {
		fmt.Fprintf(ctx.LogFile, format, args...)
	}
}

func runStage(stage *Stage, ctx *WorkflowContext) (string, error) {
	prompt := stage.Prompt
	prompt = strings.ReplaceAll(prompt, "{{.Requirement}}", ctx.Requirement)
	prompt = strings.ReplaceAll(prompt, "{{.ProjectContext}}", ctx.Results["project-context"])

	planContent, _ := os.ReadFile(filepath.Join(ctx.WorkDir, "plan.md"))
	if len(planContent) == 0 {
		planContent, _ = os.ReadFile(filepath.Join(ctx.WorkDir, "analysis.md"))
	}
	if len(planContent) == 0 {
		planContent, _ = os.ReadFile(filepath.Join(ctx.WorkDir, "refactor-plan.md"))
	}
	prompt = strings.ReplaceAll(prompt, "{{.PlanContent}}", string(planContent))

	tasksContent, _ := os.ReadFile(filepath.Join(ctx.WorkDir, "tasks.md"))
	if len(tasksContent) == 0 {
		tasksContent, _ = os.ReadFile(filepath.Join(ctx.WorkDir, "fix-tasks.md"))
	}
	if len(tasksContent) == 0 {
		tasksContent, _ = os.ReadFile(filepath.Join(ctx.WorkDir, "refactor-tasks.md"))
	}
	prompt = strings.ReplaceAll(prompt, "{{.TasksContent}}", string(tasksContent))

	reviewContent, _ := os.ReadFile(filepath.Join(ctx.WorkDir, "review.md"))
	prompt = strings.ReplaceAll(prompt, "{{.ReviewContent}}", string(reviewContent))

	verifyContent, _ := os.ReadFile(filepath.Join(ctx.WorkDir, "verify.md"))
	prompt = strings.ReplaceAll(prompt, "{{.VerifyContent}}", string(verifyContent))

	diffContent, _ := os.ReadFile(filepath.Join(ctx.WorkDir, "diff.md"))
	prompt = strings.ReplaceAll(prompt, "{{.DiffContent}}", string(diffContent))

	securityContent, _ := os.ReadFile(filepath.Join(ctx.WorkDir, "security.md"))
	prompt = strings.ReplaceAll(prompt, "{{.SecurityContent}}", string(securityContent))

	ctx.log("### Prompt\n```\n%s\n```\n\n", truncate(prompt, 1000))

	oldBackend := current
	oldModel := currentModel
	current = stage.Backend
	currentModel = stage.Model // Set model for this stage

	var result string
	if stage.Interactive {
		result = callInteractive(prompt)
	} else {
		result = call(prompt)
	}

	current = oldBackend
	currentModel = oldModel
	return result, nil
}

func getWorkflow(name string) *Workflow {
	if wf, ok := defaultWorkflows[name]; ok {
		return &wf
	}
	return nil
}

func listWorkflows() {
	fmt.Println(cyan("Available workflows:"))
	for name, wf := range defaultWorkflows {
		fmt.Printf("  %s - %s\n", green(name), wf.Name)
		for i, s := range wf.Stages {
			out := ""
			if s.OutputFile != "" {
				out = fmt.Sprintf(" → %s", s.OutputFile)
			}
			if s.Interactive {
				out += " (interactive)"
			}
			fmt.Printf("    %d. %s (%s)%s\n", i+1, s.Name, s.Backend, out)
		}
	}
}
