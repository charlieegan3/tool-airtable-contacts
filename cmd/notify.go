package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "send alerts for events related to contacts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("notify called")
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
}
