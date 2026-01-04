package main

import (
	"fmt"
	"os"

	"github.com/LarsEckart/huey/auth"
	"github.com/LarsEckart/huey/cmd"
	"github.com/LarsEckart/huey/hue"
	"github.com/LarsEckart/huey/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "huey",
	Short: "Control Philips Hue lights",
	Long:  "A CLI to control Philips Hue lights. Run without arguments for interactive TUI.",
	Run: func(command *cobra.Command, args []string) {
		cfg, err := auth.EnsureAuthenticated()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := hue.NewClient(cfg.BridgeIP, cfg.Username)
		if err := tui.Run(client); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(cmd.LightsCmd)
	rootCmd.AddCommand(cmd.LightCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
