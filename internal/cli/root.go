package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "admin-cli",
	Short: "Admin CLI for Discord Leaving Bot",
	Long:  `A management tool for Discord Leaving Bot to handle leave records, schedules and reports.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root flags can be defined here
}
