package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type WorkflowState struct {
	WorkflowName string            `json:"workflow"`
	Requirement  string            `json:"requirement"`
	CurrentStage int               `json:"currentStage"`
	Results      map[string]string `json:"results"`
	WorkDir      string            `json:"workDir"`
}

func saveCheckpoint(ctx *WorkflowContext, wfName string, stageIdx int) {
	state := WorkflowState{
		WorkflowName: wfName,
		Requirement:  ctx.Requirement,
		CurrentStage: stageIdx,
		Results:      ctx.Results,
		WorkDir:      ctx.WorkDir,
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(filepath.Join(ctx.WorkDir, "state.json"), data, 0644)
}

func loadCheckpoint(workDir string) (*WorkflowState, error) {
	data, err := os.ReadFile(filepath.Join(workDir, "state.json"))
	if err != nil {
		return nil, err
	}
	var state WorkflowState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func findLatestWorkflow() string {
	latestPath := filepath.Join(".workflow", "latest")
	if target, err := os.Readlink(latestPath); err == nil {
		return filepath.Join(".workflow", target)
	}
	return ""
}
