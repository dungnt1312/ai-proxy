package main

import (
	"os"
	"path/filepath"
	"strings"
)

func scanProjectContext() string {
	var ctx strings.Builder
	ctx.WriteString("# Project Context\n\n")

	if _, err := os.Stat("go.mod"); err == nil {
		ctx.WriteString("## Tech Stack: Go\n")
		if content, _ := os.ReadFile("go.mod"); len(content) > 0 {
			ctx.WriteString("```\n" + string(content) + "```\n\n")
		}
	} else if _, err := os.Stat("package.json"); err == nil {
		ctx.WriteString("## Tech Stack: Node.js\n")
		if content, _ := os.ReadFile("package.json"); len(content) > 0 {
			lines := strings.Split(string(content), "\n")
			if len(lines) > 20 {
				lines = lines[:20]
			}
			ctx.WriteString("```json\n" + strings.Join(lines, "\n") + "\n```\n\n")
		}
	} else if _, err := os.Stat("requirements.txt"); err == nil {
		ctx.WriteString("## Tech Stack: Python\n")
	} else if _, err := os.Stat("Cargo.toml"); err == nil {
		ctx.WriteString("## Tech Stack: Rust\n")
	}

	for _, readme := range []string{"README.md", "readme.md", "README"} {
		if content, err := os.ReadFile(readme); err == nil {
			ctx.WriteString("## README\n")
			lines := strings.Split(string(content), "\n")
			if len(lines) > 30 {
				lines = lines[:30]
				lines = append(lines, "...")
			}
			ctx.WriteString(strings.Join(lines, "\n") + "\n\n")
			break
		}
	}

	ctx.WriteString("## Project Structure\n```\n")
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || strings.HasPrefix(path, ".git") || strings.HasPrefix(path, ".workflow") ||
			strings.HasPrefix(path, "node_modules") || strings.HasPrefix(path, "vendor") {
			return filepath.SkipDir
		}
		if info.IsDir() && path != "." {
			return nil
		}
		if !info.IsDir() {
			ctx.WriteString(path + "\n")
		}
		return nil
	})
	ctx.WriteString("```\n")

	return ctx.String()
}
