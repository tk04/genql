package cmd

import (
	"errors"
	"fmt"
	"genql/prismaUtil"
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
		addTypes(args[0], model)
		addResolvers(args[0], model)

	},
}

func addTypes(modelName string, model prismaUtil.Model) {
	resolverPath := "./src/resolvers/" + modelName

	filePath := resolverPath + "/types.ts"
	if checkFileExists(filePath) {
		fmt.Printf("file (%s) already exists\n", filePath)
		os.Exit(1)
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	header := "import { Field, InputType, ObjectType } from \"type-graphql\"\n\n"
	f.WriteString(header + model.ObjectType() + "\n" + model.CreateInputType() + "\n" + model.UpdateInputType())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addResolvers(modelName string, model prismaUtil.Model) {
	resolverPath := "./src/resolvers/" + modelName

	filePath := resolverPath + "/index.ts"
	if checkFileExists(filePath) {
		fmt.Printf("file (%s) already exists\n", filePath)
		os.Exit(1)
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// TODO: write resolvers to filePath
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}
	fmt.Printf("error ecountered find file @ path: %s\n", filePath)
	os.Exit(1)

	return false
}
