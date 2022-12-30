package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genql",
	Short: "genql is a server side GraphQL & Prisma code generator",
	Long:  "A code generator that reliably generates Prisma database schemas followed by CRUD GraphQL resolvers for each generated model.",
}

func Execute() {
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(modelCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
