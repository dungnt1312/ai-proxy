package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagBackend string
	flagList    bool
	flagInit    bool
)

var rootCmd = &cobra.Command{
	Use:   "proxy [prompt]",
	Short: "AI CLI Proxy - unified interface for multiple AI CLIs",
	Run: func(cmd *cobra.Command, args []string) {
		config = loadConfig()

		if flagInit {
			if err := initProject(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if flagList {
			fmt.Println("Available backends:")
			for k, v := range config.Backends {
				mark := "  "
				if k == config.Default {
					mark = "* "
				}
				fmt.Printf("%s%s (%s) - %s\n", mark, k, v.Name, v.Cmd)
			}
			return
		}

		if flagBackend != "" {
			if _, ok := config.Backends[flagBackend]; ok {
				current = flagBackend
			} else {
				fmt.Printf("Unknown backend: %s\n", flagBackend)
				os.Exit(1)
			}
		} else {
			current = config.Default
		}

		if len(args) > 0 {
			prompt := args[0]
			call(prompt)
			return
		}

		runInteractive()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&flagBackend, "backend", "b", "", "Backend to use (claude, kiro, gemini, cursor)")
	rootCmd.Flags().BoolVarP(&flagList, "list", "l", false, "List available backends")
	rootCmd.Flags().BoolVar(&flagInit, "init", false, "Initialize project config (.ai-proxy/config.json)")
}

// Execute runs the root CLI command.
//
// This is the entrypoint used by main(). It parses flags/args and either:
//   - initializes project config (--init),
//   - lists configured backends (--list), or
//   - starts interactive mode (no args) / sends a one-shot prompt (args present).
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
