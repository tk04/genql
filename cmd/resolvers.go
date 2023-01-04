package cmd

import (
	"fmt"
	"github.com/tk04/genql/prismaUtil"
	"github.com/tk04/genql/resolvers"
	"os"

	"github.com/spf13/cobra"
)

var resolversCmd = &cobra.Command{
	Use:   "resolvers",
	Short: "Generate GraphQL resolvers for a Prisma Model",
	Long:  "Generate CRUD GraphQL resolvers for a Prisma Model.\n\n Usage: genql resolvers [model name].",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg, _ := cmd.Flags().GetStringArray("Except")

		val := struct{}{}
		funcs := map[string]struct{}{
			"get":    val,
			"create": val,
			"update": val,
			"delete": val,
		}
		for _, val := range arg {
			delete(funcs, val)
		}
		include := []string{}
		for k := range funcs {
			include = append(include, k)
		}

		resolverPath := "./src/resolvers/" + args[0]
		err := os.MkdirAll(resolverPath, os.ModePerm)
		if err != nil {
			fmt.Println("err")
			os.Exit(1)
		}

		model := prismaUtil.GetModel(args[0])
		resolver := resolvers.Resolver{Model: model, Functions: include}
		resolver.CreateFiles()
	},
}
