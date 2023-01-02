package cmd

import (
	"fmt"
	"genql/prismaUtil"
	"genql/resolvers"
	"os"

	"github.com/spf13/cobra"
)

var resolversCmd = &cobra.Command{
	Use:   "resolvers",
	Short: "Generate GraphQL resolvers for a Prisma Model",
	Long:  "Generate CRUD GraphQL resolvers for a Prisma Model.\n\n Usage: genql resolvers [model name].",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resolverPath := "./src/resolvers/" + args[0]
		err := os.MkdirAll(resolverPath, os.ModePerm)
		if err != nil {
			fmt.Println("err")
			os.Exit(1)
		}

		model := prismaUtil.GetModel(args[0])
		resolvers.CreateResolver(model)
	},
}
