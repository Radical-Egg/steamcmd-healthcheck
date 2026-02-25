package cmd

import (
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:           "program",
	Short:         "Send a UDP heartbeat to steamcmd gameserver",
	SilenceUsage:  true, // don't print usage on runtime errors
	SilenceErrors: true, // don't let Cobra print errors automatically
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging to stderr")
}
