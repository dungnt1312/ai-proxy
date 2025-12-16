package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type BackendConfig struct {
	Name       string   `json:"name"`
	Cmd        string   `json:"cmd"`
	Args       []string `json:"args"`
	PromptFlag string   `json:"promptFlag"`
	ResumeFlag string   `json:"resumeFlag"`
	ModelFlag  string   `json:"modelFlag,omitempty"` // e.g., "--model" for claude/kiro
}

type Config struct {
	Default   string                   `json:"default"`
	Backends  map[string]BackendConfig `json:"backends"`
	Workflows map[string]Workflow      `json:"workflows,omitempty"`
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
				Name:       "Cursor",
				Cmd:        "cursor",
				Args:       []string{},
				PromptFlag: "",
				ResumeFlag: "",
				ModelFlag:  "--model",
			},
		},
	}
}
