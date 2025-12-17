package main

import (
	"fmt"

	"github.com/fatih/color"
)

func init() {
	// Ensure examples have stable output (no ANSI escape sequences).
	color.NoColor = true
}

func ExampleWorkflow_DryRun() {
	wf := &Workflow{
		Name: "Demo",
		Stages: []Stage{
			{Name: "plan", Backend: "gemini", OutputFile: "plan.md"},
			{Name: "execute", Backend: "claude", Interactive: true, Skippable: true},
		},
	}

	_ = wf.DryRun("add feature X")

	// Output:
	//
	// ▶ DRY RUN: Demo
	// │ Requirement: add feature X
	//
	// │ Stage 1: plan (gemini)
	// │   → plan.md
	// │ Stage 2: execute (claude) (interactive) [will ask to skip]
	//
	// ! This is a dry run. No changes will be made.
	// │ Run without --dry-run to execute.
}

func ExampleSkill_ToStage() {
	s := &Skill{
		Name: "code-review",
		Stage: SkillStage{
			Backend:     "kiro",
			Interactive: false,
			OutputFile:  "review.md",
		},
		Prompt: "Review: {{.DiffContent}}",
	}

	stage := s.ToStage()
	fmt.Printf("%s %s %t %s\n", stage.Name, stage.Backend, stage.Interactive, stage.OutputFile)

	// Output:
	// code-review kiro false review.md
}
