package cmd

import (
	"github.com/spf13/cobra"
)

var resolversCmd = &cobra.Command{
	Use:   "resolvers",
	Short: "Generate GraphQL resolvers for a Prisma Model",
	Long:  "Generate CRUD GraphQL resolvers for a Prisma Model.\n\n Usage: genql resolvers [model name].",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
