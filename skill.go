package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Skill struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Version     string       `yaml:"version"`
	Author      string       `yaml:"author"`
	Stage       SkillStage   `yaml:"stage"`
	Inputs      []SkillInput `yaml:"inputs"`
	Tags        []string     `yaml:"tags"`
	Prompt      string       `yaml:"-"` // Loaded from prompt.md
	Path        string       `yaml:"-"` // Skill directory path
}

type SkillStage struct {
	Backend     string `yaml:"backend"`
	Model       string `yaml:"model"`
	Interactive bool   `yaml:"interactive"`
	OutputFile  string `yaml:"outputFile"`
}

type SkillInput struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default"`
}

var skills = make(map[string]*Skill)

func getSkillsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ai-proxy", "skills")
}

func loadSkills() {
	// Load from global ~/.ai-proxy/skills/
	loadSkillsFromDir(getSkillsDir())

	// Load from project .ai-proxy/skills/ (override global)
	loadSkillsFromDir(filepath.Join(".ai-proxy", "skills"))
}

func loadSkillsFromDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillPath := filepath.Join(dir, e.Name())
		skill, err := loadSkill(skillPath)
		if err != nil {
			continue
		}
		skills[skill.Name] = skill
	}
}

func loadSkill(path string) (*Skill, error) {
	// Load skill.yaml
	yamlPath := filepath.Join(path, "skill.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	var skill Skill
	if err := yaml.Unmarshal(data, &skill); err != nil {
		return nil, err
	}

	// Load prompt.md
	promptPath := filepath.Join(path, "prompt.md")
	promptData, err := os.ReadFile(promptPath)
	if err != nil {
		return nil, err
	}
	skill.Prompt = string(promptData)
	skill.Path = path

	return &skill, nil
}

func (s *Skill) Run(inputs map[string]string, ctx *WorkflowContext) (string, error) {
	// Build prompt from template
	prompt := s.Prompt

	// Replace inputs
	for _, input := range s.Inputs {
		placeholder := fmt.Sprintf("{{.%s}}", input.Name)
		value := inputs[input.Name]
		if value == "" {
			value = input.Default
		}
		if value == "" && input.Required {
			return "", fmt.Errorf("missing required input: %s", input.Name)
		}
		prompt = strings.ReplaceAll(prompt, placeholder, value)
	}

	// Replace standard context variables if ctx provided
	if ctx != nil {
		prompt = strings.ReplaceAll(prompt, "{{.ProjectContext}}", ctx.Results["project-context"])
		prompt = strings.ReplaceAll(prompt, "{{.Requirement}}", ctx.Requirement)
		prompt = strings.ReplaceAll(prompt, "{{.DiffContent}}", ctx.Results["diff"])
	}

	// Execute
	oldBackend := current
	oldModel := currentModel

	if s.Stage.Backend != "" {
		current = s.Stage.Backend
	}
	if s.Stage.Model != "" {
		currentModel = s.Stage.Model
	}

	var result string
	if s.Stage.Interactive {
		result = callInteractive(prompt)
	} else {
		result = call(prompt)
	}

	current = oldBackend
	currentModel = oldModel

	return result, nil
}

func (s *Skill) ToStage() Stage {
	return Stage{
		Name:        s.Name,
		Backend:     s.Stage.Backend,
		Model:       s.Stage.Model,
		Prompt:      s.Prompt,
		OutputFile:  s.Stage.OutputFile,
		Interactive: s.Stage.Interactive,
	}
}

func listSkills() {
	if len(skills) == 0 {
		fmt.Println(dim("No skills installed"))
		fmt.Println(dim("Add skills to ~/.ai-proxy/skills/ or .ai-proxy/skills/"))
		fmt.Println(dim("Or install: /skill install <github-url>"))
		return
	}

	fmt.Println(cyan("Available skills:"))
	for name, skill := range skills {
		tags := ""
		if len(skill.Tags) > 0 {
			tags = dim(fmt.Sprintf(" [%s]", strings.Join(skill.Tags, ", ")))
		}
		fmt.Printf("  %s - %s%s\n", green(name), skill.Description, tags)
	}
}

func getSkill(name string) *Skill {
	return skills[name]
}

func runSkillCommand(args []string) {
	if len(args) == 0 {
		listSkills()
		return
	}

	switch args[0] {
	case "install":
		if len(args) < 2 {
			fmt.Println("Usage: /skill install <github-url>")
			fmt.Println("Example: /skill install https://github.com/user/repo/tree/main/skills/my-skill")
			return
		}
		installSkill(args[1])
		return

	case "remove":
		if len(args) < 2 {
			fmt.Println("Usage: /skill remove <name>")
			return
		}
		removeSkill(args[1])
		return

	case "info":
		if len(args) < 2 {
			fmt.Println("Usage: /skill info <name>")
			return
		}
		showSkillInfo(args[1])
		return
	}

	// Run skill
	skillName := args[0]
	skill := getSkill(skillName)
	if skill == nil {
		fmt.Printf("%s Skill not found: %s\n", red("!"), skillName)
		listSkills()
		return
	}

	// Parse inputs from args: --input=value
	inputs := make(map[string]string)
	var overrideBackend string
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
			if len(parts) == 2 {
				if parts[0] == "backend" {
					overrideBackend = parts[1]
				} else {
					inputs[parts[0]] = parts[1]
				}
			}
		}
	}

	// Override backend if specified
	if overrideBackend != "" {
		skill.Stage.Backend = overrideBackend
	}

	// Check required inputs
	for _, input := range skill.Inputs {
		if input.Required && inputs[input.Name] == "" {
			fmt.Printf("%s Missing required input: --%s\n", red("!"), input.Name)
			fmt.Printf("Usage: /skill %s", skillName)
			for _, i := range skill.Inputs {
				if i.Required {
					fmt.Printf(" --%s=<value>", i.Name)
				} else {
					fmt.Printf(" [--%s=<value>]", i.Name)
				}
			}
			fmt.Println(" [--backend=<name>]")
			return
		}
	}

	backendInfo := skill.Stage.Backend
	if backendInfo == "" {
		backendInfo = current
	}
	fmt.Printf("%s Running skill: %s (%s)\n", cyan("▶"), skill.Name, backendInfo)
	result, err := skill.Run(inputs, nil)
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	if skill.Stage.OutputFile != "" {
		os.WriteFile(skill.Stage.OutputFile, []byte(result), 0644)
		fmt.Printf("%s Output saved to: %s\n", green("✓"), skill.Stage.OutputFile)
	}
}

func installSkill(url string) {
	// Convert GitHub URL to raw URLs
	// https://github.com/user/repo/tree/main/skills/my-skill
	// -> https://raw.githubusercontent.com/user/repo/main/skills/my-skill/skill.yaml

	var rawBase string
	if strings.Contains(url, "github.com") {
		// Parse GitHub URL
		url = strings.TrimSuffix(url, "/")
		parts := strings.Split(url, "/")
		if len(parts) < 7 {
			fmt.Printf("%s Invalid GitHub URL\n", red("!"))
			return
		}
		user := parts[3]
		repo := parts[4]
		branch := parts[6]
		path := strings.Join(parts[7:], "/")
		rawBase = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", user, repo, branch, path)
	} else {
		rawBase = url
	}

	// Download skill.yaml
	fmt.Printf("%s Downloading skill from %s...\n", cyan("↓"), url)

	yamlURL := rawBase + "/skill.yaml"
	yamlData, err := downloadFile(yamlURL)
	if err != nil {
		fmt.Printf("%s Failed to download skill.yaml: %v\n", red("!"), err)
		return
	}

	var skill Skill
	if err := yaml.Unmarshal(yamlData, &skill); err != nil {
		fmt.Printf("%s Invalid skill.yaml: %v\n", red("!"), err)
		return
	}

	// Download prompt.md
	promptURL := rawBase + "/prompt.md"
	promptData, err := downloadFile(promptURL)
	if err != nil {
		fmt.Printf("%s Failed to download prompt.md: %v\n", red("!"), err)
		return
	}

	// Create skill directory
	skillDir := filepath.Join(getSkillsDir(), skill.Name)
	os.MkdirAll(skillDir, 0755)

	// Save files
	os.WriteFile(filepath.Join(skillDir, "skill.yaml"), yamlData, 0644)
	os.WriteFile(filepath.Join(skillDir, "prompt.md"), promptData, 0644)

	// Reload skills
	skill.Path = skillDir
	skill.Prompt = string(promptData)
	skills[skill.Name] = &skill

	fmt.Printf("%s Installed skill: %s v%s\n", green("✓"), skill.Name, skill.Version)
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func removeSkill(name string) {
	skill := getSkill(name)
	if skill == nil {
		fmt.Printf("%s Skill not found: %s\n", red("!"), name)
		return
	}

	// Only remove from global skills dir
	globalPath := filepath.Join(getSkillsDir(), name)
	if skill.Path != globalPath {
		fmt.Printf("%s Cannot remove project-local skill\n", yellow("!"))
		return
	}

	if err := os.RemoveAll(globalPath); err != nil {
		fmt.Printf("%s Failed to remove: %v\n", red("!"), err)
		return
	}

	delete(skills, name)
	fmt.Printf("%s Removed skill: %s\n", green("✓"), name)
}

func showSkillInfo(name string) {
	skill := getSkill(name)
	if skill == nil {
		fmt.Printf("%s Skill not found: %s\n", red("!"), name)
		return
	}

	fmt.Printf("%s %s\n", cyan("Skill:"), skill.Name)
	fmt.Printf("%s %s\n", dim("Description:"), skill.Description)
	fmt.Printf("%s %s\n", dim("Version:"), skill.Version)
	fmt.Printf("%s %s\n", dim("Author:"), skill.Author)
	fmt.Printf("%s %s\n", dim("Backend:"), skill.Stage.Backend)
	fmt.Printf("%s %s\n", dim("Path:"), skill.Path)

	if len(skill.Inputs) > 0 {
		fmt.Printf("%s\n", dim("Inputs:"))
		for _, input := range skill.Inputs {
			req := ""
			if input.Required {
				req = red("*")
			}
			fmt.Printf("  --%s%s - %s\n", input.Name, req, input.Description)
		}
	}
}
