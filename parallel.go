package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func checkCondition(cond string, ctx *WorkflowContext) bool {
	if cond == "" {
		return true
	}
	switch {
	case strings.HasPrefix(cond, "file:"):
		file := strings.TrimPrefix(cond, "file:")
		_, err := os.Stat(file)
		return err == nil
	case strings.HasPrefix(cond, "!file:"):
		file := strings.TrimPrefix(cond, "!file:")
		_, err := os.Stat(file)
		return err != nil
	case strings.HasPrefix(cond, "has:"):
		ext := strings.TrimPrefix(cond, "has:")
		found := false
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ext) {
				found = true
			}
			return nil
		})
		return found
	case cond == "go":
		_, err := os.Stat("go.mod")
		return err == nil
	case cond == "node":
		_, err := os.Stat("package.json")
		return err == nil
	case cond == "docker":
		_, err := os.Stat("Dockerfile")
		return err == nil
	}
	return true
}

type ParallelResult struct {
	Name   string
	Result string
	Err    error
}

func runParallelStages(stages []Stage, ctx *WorkflowContext) []ParallelResult {
	var wg sync.WaitGroup
	results := make([]ParallelResult, len(stages))

	for i, stage := range stages {
		wg.Add(1)
		go func(idx int, s Stage) {
			defer wg.Done()
			result, err := runStage(&s, ctx)
			results[idx] = ParallelResult{Name: s.Name, Result: result, Err: err}
		}(i, stage)
	}

	wg.Wait()
	return results
}
