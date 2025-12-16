package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Message struct {
	Role    string
	Content string
}

var (
	config  *Config
	current string
	history = []Message{}

	cyan   = color.New(color.FgCyan).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	dim    = color.New(color.Faint).SprintFunc()
)

var currentModel string // Model for current stage

func buildArgs(prompt string) []string {
	b := config.Backends[current]
	args := append([]string{}, b.Args...)

	if b.PromptFlag != "" {
		args = append(args, b.PromptFlag, prompt)
	} else {
		args = append(args, prompt)
	}

	if len(history) > 0 && b.ResumeFlag != "" {
		args = append(args, b.ResumeFlag)
	}

	if currentModel != "" && b.ModelFlag != "" {
		args = append(args, b.ModelFlag, currentModel)
	}

	if current == "kiro" {
		args = append(args, "--no-interactive")
	}

	return args
}

func call(prompt string) string {
	b := config.Backends[current]
	args := buildArgs(prompt)

	fmt.Printf("%s %s %s\n", dim("→"), dim(b.Cmd), dim(truncate(strings.Join(args, " "), 80)))

	start := time.Now()
	cmd := exec.Command(b.Cmd, args...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	var response strings.Builder
	go io.Copy(os.Stderr, stderr)

	buf := make([]byte, 256)
	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			os.Stdout.Write(buf[:n])
			response.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	cmd.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\n%s\n", dim(fmt.Sprintf("(%s)", elapsed.Round(time.Millisecond))))

	return strings.TrimSpace(response.String())
}

func callInteractive(prompt string) string {
	b := config.Backends[current]

	var args []string
	if current == "claude" {
		args = []string{prompt}
	} else {
		args = buildArgs(prompt)
	}

	fmt.Printf("%s %s %s %s\n", dim("→"), dim(b.Cmd), dim(truncate(strings.Join(args, " "), 60)), yellow("(interactive)"))
	fmt.Printf("%s Press Ctrl+C when done\n\n", dim("│"))

	cmd := exec.Command(b.Cmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
	return "(interactive session completed)"
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

func handleCommand(input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	switch parts[0] {
	case "/switch", "/s":
		if len(parts) < 2 {
			fmt.Println("Usage: /switch <backend>")
			return true
		}
		if _, ok := config.Backends[parts[1]]; ok {
			current = parts[1]
			config.Default = current
			saveConfig(config)
			fmt.Printf("%s Switched to %s\n", green("✓"), config.Backends[current].Name)
		} else {
			fmt.Printf("%s Unknown: %s\n", yellow("!"), parts[1])
		}
		return true

	case "/list", "/l":
		fmt.Println(cyan("Backends:"))
		for k, v := range config.Backends {
			if k == current {
				fmt.Printf("  %s %s (%s)\n", green("*"), green(k), v.Name)
			} else {
				fmt.Printf("    %s (%s)\n", k, v.Name)
			}
		}
		return true

	case "/clear", "/c":
		history = []Message{}
		fmt.Printf("%s Cleared\n", green("✓"))
		return true

	case "/history":
		if len(history) == 0 {
			fmt.Println(dim("Empty"))
			return true
		}
		for i, m := range history {
			c := m.Content
			if len(c) > 50 {
				c = c[:50] + "..."
			}
			role := cyan(m.Role)
			if m.Role == "user" {
				role = green(m.Role)
			}
			fmt.Printf("%s %d. [%s] %s\n", dim("│"), i+1, role, c)
		}
		return true

	case "/config":
		fmt.Printf("Config: %s\n", cyan(configPath))
		return true

	case "/init":
		if err := initProject(); err != nil {
			fmt.Printf("%s %v\n", red("Error:"), err)
		}
		return true

	case "/resume":
		if err := resumeWorkflow(); err != nil {
			fmt.Printf("%s %v\n", red("Error:"), err)
		}
		return true

	case "/workflow", "/w":
		if len(parts) < 2 {
			listWorkflows()
			return true
		}
		if parts[1] == "history" {
			showWorkflowHistory()
			return true
		}
		wfName := parts[1]
		if strings.HasPrefix(wfName, "--dry-run") && len(parts) > 2 {
			dryRun = true
			wfName = parts[2]
			parts = append(parts[:1], parts[2:]...)
		}
		wf := getWorkflow(wfName)
		if wf == nil {
			fmt.Printf("%s Unknown workflow: %s\n", yellow("!"), wfName)
			listWorkflows()
			return true
		}
		var req string
		if len(parts) > 2 {
			req = strings.Join(parts[2:], " ")
		} else {
			fmt.Print("Requirement: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			req = scanner.Text()
		}
		if err := wf.Run(req); err != nil {
			fmt.Printf("%s %v\n", red("Error:"), err)
		}
		dryRun = false
		return true

	case "/help", "/?":
		fmt.Println(cyan("Commands:"))
		fmt.Println("  /init                - Init project config")
		fmt.Println("  /switch <name>       - Switch backend")
		fmt.Println("  /list                - List backends")
		fmt.Println("  /workflow <name>     - Run workflow")
		fmt.Println("  /workflow history    - Show workflow history")
		fmt.Println("  /workflow --dry-run <name> - Preview workflow")
		fmt.Println("  /resume              - Resume last workflow")
		fmt.Println("  /clear               - Clear history")
		fmt.Println("  /config              - Show config path")
		fmt.Println("  quit                 - Exit")
		return true
	}
	return false
}

func runInteractive() {
	loadProjectConfig()

	fmt.Println(green("🔀 AI Proxy CLI"))
	fmt.Printf("Backend: %s %s\n\n", cyan(config.Backends[current].Name), dim("(/? for help)"))

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s ", green(fmt.Sprintf("[%s]>", current)))
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "quit" {
			fmt.Println(yellow("Bye!"))
			break
		}
		if strings.HasPrefix(input, "/") {
			handleCommand(input)
			continue
		}

		history = append(history, Message{"user", input})
		resp := call(input)
		if resp != "" {
			history = append(history, Message{"assistant", resp})
		}
	}
}

func main() {
	Execute()
}
