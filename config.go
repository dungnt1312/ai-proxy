package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// BackendConfig describes how to execute a backend CLI (command, args, and flag conventions).
type BackendConfig struct {
	// Name is a human-friendly backend label shown in the UI (e.g. "Claude", "Gemini").
	Name string `json:"name"`
	// Cmd is the executable name to invoke (must be available on PATH).
	Cmd string `json:"cmd"`
	// Args are default arguments to pass to the backend CLI.
	Args []string `json:"args"`
	// PromptFlag is the CLI flag that accepts the prompt text (e.g. "-p").
	// If empty, the prompt is appended as a positional argument.
	PromptFlag string `json:"promptFlag"`
	// ResumeFlag is the CLI flag to resume/continue a prior chat session, if supported.
	ResumeFlag string `json:"resumeFlag"`
	// ModelFlag is the CLI flag that selects a model (e.g. "--model" or "-m").
	ModelFlag string `json:"modelFlag,omitempty"` // e.g., "--model" for claude/kiro
}

// Config is the global configuration loaded from `~/.ai-proxy.json`.
// It defines available backends and optional workflow overrides.
type Config struct {
	// Default is the backend key to use when none is explicitly selected.
	Default string `json:"default"`
	// Backends maps a backend key (e.g. "claude") to its execution configuration.
	Backends map[string]BackendConfig `json:"backends"`
	// Workflows optionally overrides/extends built-in workflows with project-defined workflows.
	Workflows map[string]Workflow `json:"workflows,omitempty"`
}

var configPath string

func init() {
	home, _ := os.UserHomeDir()
	configPath = filepath.Join(home, ".ai-proxy.json")
}

func loadConfig() *Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig()
	}
	var cfg Config
	if json.Unmarshal(data, &cfg) != nil {
		return defaultConfig()
	}
	return &cfg
}

func saveConfig(cfg *Config) error {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(configPath, data, 0644)
}

func defaultConfig() *Config {
	return &Config{
		Default: "claude",
		Backends: map[string]BackendConfig{
			"claude": {
				Name:       "Claude",
				Cmd:        "claude",
				Args:       []string{},
				PromptFlag: "-p",
				ResumeFlag: "--continue",
				ModelFlag:  "--model",
			},
			"kiro": {
				Name:       "Kiro",
				Cmd:        "kiro-cli",
				Args:       []string{"chat"},
				PromptFlag: "",
				ResumeFlag: "--resume",
				ModelFlag:  "--model",
			},
			"gemini": {
				Name:       "Gemini",
				Cmd:        "gemini",
				Args:       []string{},
				PromptFlag: "",
				ResumeFlag: "--resume",
				ModelFlag:  "-m",
			},
			"cursor": {
				Name: "Cursor",
				// Cursor's agent binary is typically named "cursor-agent" (as documented in README.md).
				Cmd:        "cursor-agent",
				Args:       []string{},
				PromptFlag: "",
				ResumeFlag: "",
				ModelFlag:  "--model",
			},
		},
	}
}
