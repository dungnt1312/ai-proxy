package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// VerifyResult is the outcome of running repository verification checks (build/vet/test, etc.).
// Output is formatted as Markdown for inclusion in workflow artifacts.
type VerifyResult struct {
	Passed bool
	Output string
}

func autoVerify() VerifyResult {
	var results strings.Builder
	allPassed := true

	if _, err := os.Stat("go.mod"); err == nil {
		results.WriteString("## Go Checks\n\n")

		results.WriteString("### Build\n")
		if out, err := exec.Command("go", "build", "./...").CombinedOutput(); err != nil {
			results.WriteString("❌ FAILED\n```\n" + string(out) + "```\n\n")
			allPassed = false
		} else {
			results.WriteString("✅ PASSED\n\n")
		}

		results.WriteString("### Vet\n")
		if out, err := exec.Command("go", "vet", "./...").CombinedOutput(); err != nil {
			results.WriteString("⚠️ ISSUES\n```\n" + string(out) + "```\n\n")
		} else {
			results.WriteString("✅ PASSED\n\n")
		}

		results.WriteString("### Tests\n")
		out, err := exec.Command("go", "test", "./...", "-v").CombinedOutput()
		if err != nil {
			results.WriteString("❌ FAILED\n```\n" + string(out) + "```\n\n")
			allPassed = false
		} else if strings.Contains(string(out), "no test files") {
			results.WriteString("⚠️ No tests found\n\n")
		} else {
			results.WriteString("✅ PASSED\n```\n" + truncate(string(out), 500) + "```\n\n")
		}
	}

	if _, err := os.Stat("package.json"); err == nil {
		results.WriteString("## Node.js Checks\n\n")

		results.WriteString("### Tests\n")
		if out, err := exec.Command("npm", "test").CombinedOutput(); err != nil {
			results.WriteString("❌ FAILED\n```\n" + string(out) + "```\n\n")
			allPassed = false
		} else {
			results.WriteString("✅ PASSED\n\n")
		}
	}

	return VerifyResult{Passed: allPassed, Output: results.String()}
}

func runVerifyStage(workDir string) (string, bool) {
	fmt.Printf("%s Running auto-verify...\n", dim("│"))

	result := autoVerify()

	status := "✅ ALL PASSED"
	if !result.Passed {
		status = "❌ SOME CHECKS FAILED"
	}

	output := fmt.Sprintf("# Auto-Verify Results\n\n**Status:** %s\n\n%s", status, result.Output)
	return output, result.Passed
}
