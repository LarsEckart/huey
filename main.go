package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "huey",
	Short: "Control Philips Hue lights",
	Long:  "A CLI to control Philips Hue lights. Run without arguments for interactive TUI.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Launch TUI when no subcommand given
		fmt.Println("huey - Philips Hue CLI (TUI coming soon)")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
