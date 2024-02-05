package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "mosaic",
		Short:         "Generate mosaic videos from many inputs",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(Work())
	rootCmd.AddCommand(Store())
	rootCmd.AddCommand(Player())
	rootCmd.AddCommand(Simulate())

	return rootCmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
