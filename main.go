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

func rootAction(command *cobra.Command, args []string) error {
	cfg, err := auth.EnsureAuthenticated()
	if err != nil {
		return fmt.Errorf("ensure authentication: %w", err)
	}

	client := hue.NewClient(cfg.BridgeIP, cfg.Username)
	if err := tui.Run(client); err != nil {
		return fmt.Errorf("run tui: %w", err)
	}

	return nil
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "huey",
		Short:         "Control Philips Hue lights",
		Long:          "A CLI to control Philips Hue lights. Run without arguments for interactive TUI.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       appVersion(),
		RunE:          rootAction,
	}
	rootCmd.SetVersionTemplate("{{.Name}} version {{.Version}}\n")

	rootCmd.AddCommand(cmd.LightsCmd)
	rootCmd.AddCommand(cmd.LightCmd)
	rootCmd.AddCommand(cmd.GroupsCmd)
	rootCmd.AddCommand(cmd.GroupCmd)
	rootCmd.AddCommand(cmd.GroupCreateCmd)
	rootCmd.AddCommand(cmd.ScenesCmd)
	rootCmd.AddCommand(cmd.SceneCmd)
	rootCmd.AddCommand(cmd.SceneCreateCmd)

	return rootCmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
