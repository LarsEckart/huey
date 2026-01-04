package main

import (
	"fmt"
	"os"

	"github.com/lars/huey/auth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "huey",
	Short: "Control Philips Hue lights",
	Long:  "A CLI to control Philips Hue lights. Run without arguments for interactive TUI.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// TODO: Launch TUI when no subcommand given
		fmt.Printf("huey - Connected to bridge at %s\n", cfg.BridgeIP)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
