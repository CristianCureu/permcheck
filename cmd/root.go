package cmd

import (
	"os"

	"github.com/cristiancureu/permcheck/internal"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "permcheck",
	Short: "A CLI tool to audit file permissions",
	Long:  "permcheck recursively scans a directory and flags files with insecure permissions.",
}

func Execute() error {
	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))

	cfg := &internal.Config{
		IsTerminal:    isTerminal,
		ColorsEnabled: isTerminal,
	}

	rootCmd.AddCommand(NewScanCmd(cfg))

	return rootCmd.Execute()
}
