package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const localConfigDir = ".ai-proxy"
const localConfigFile = ".ai-proxy/config.json"

type ProjectConfig struct {
	Workflows map[string]Workflow `json:"workflows"`
}

func initProject() error {
	if err := os.MkdirAll(localConfigDir, 0755); err != nil {
		return err
	}

	cfg := ProjectConfig{
		Workflows: map[string]Workflow{
			"feature": defaultWorkflows["feature"],
			"bugfix":  defaultWorkflows["bugfix"],
		},
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(localConfigFile, data, 0644); err != nil {
		return err
	}

	fmt.Printf("%s Created %s\n", green("✓"), localConfigFile)
	fmt.Println(dim("Edit this file to customize workflows for this project"))
	return nil
}

func loadProjectConfig() {
	if data, err := os.ReadFile(localConfigFile); err == nil {
		var cfg ProjectConfig
		if json.Unmarshal(data, &cfg) == nil {
			for name, wf := range cfg.Workflows {
				defaultWorkflows[name] = wf
			}
			fmt.Printf("%s Loaded project config\n", dim("●"))
			return
		}
	}

	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		configPath := filepath.Join(dir, localConfigFile)
		if data, err := os.ReadFile(configPath); err == nil {
			var cfg ProjectConfig
			if json.Unmarshal(data, &cfg) == nil {
				for name, wf := range cfg.Workflows {
					defaultWorkflows[name] = wf
				}
				fmt.Printf("%s Loaded config from %s\n", dim("●"), configPath)
				return
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}
