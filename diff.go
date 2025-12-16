package main

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FileSnapshot struct {
	Files map[string]string // path -> md5 hash
}

func takeSnapshot() *FileSnapshot {
	snap := &FileSnapshot{Files: make(map[string]string)}

	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasPrefix(path, ".") || strings.Contains(path, "node_modules") ||
			strings.Contains(path, "vendor") || strings.Contains(path, ".workflow") {
			return nil
		}
		if content, err := os.ReadFile(path); err == nil {
			snap.Files[path] = fmt.Sprintf("%x", md5.Sum(content))
		}
		return nil
	})

	return snap
}

func (before *FileSnapshot) Diff(after *FileSnapshot) string {
	var diff strings.Builder
	diff.WriteString("# Changes Made\n\n")

	var newFiles []string
	for path := range after.Files {
		if _, exists := before.Files[path]; !exists {
			newFiles = append(newFiles, path)
		}
	}
	if len(newFiles) > 0 {
		diff.WriteString("## New Files\n")
		for _, f := range newFiles {
			diff.WriteString(fmt.Sprintf("- `%s`\n", f))
		}
		diff.WriteString("\n")
	}

	var modifiedFiles []string
	for path, hash := range after.Files {
		if beforeHash, exists := before.Files[path]; exists && beforeHash != hash {
			modifiedFiles = append(modifiedFiles, path)
		}
	}
	if len(modifiedFiles) > 0 {
		diff.WriteString("## Modified Files\n")
		for _, f := range modifiedFiles {
			diff.WriteString(fmt.Sprintf("- `%s`\n", f))
		}
		diff.WriteString("\n")
	}

	var deletedFiles []string
	for path := range before.Files {
		if _, exists := after.Files[path]; !exists {
			deletedFiles = append(deletedFiles, path)
		}
	}
	if len(deletedFiles) > 0 {
		diff.WriteString("## Deleted Files\n")
		for _, f := range deletedFiles {
			diff.WriteString(fmt.Sprintf("- `%s`\n", f))
		}
		diff.WriteString("\n")
	}

	diff.WriteString("## File Contents\n\n")
	for _, f := range append(newFiles, modifiedFiles...) {
		content, _ := os.ReadFile(f)
		diff.WriteString(fmt.Sprintf("### %s\n```\n%s\n```\n\n", f, truncate(string(content), 1000)))
	}

	return diff.String()
}
