package cmd

import (
	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "send alerts for events related to contacts",
}

func init() {
	rootCmd.AddCommand(notifyCmd)
}
