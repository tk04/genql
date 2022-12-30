package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test cmd",
	Long:  "test cmd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test cmd")
		fmt.Println(args)
	},
}
