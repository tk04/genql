package resolvers

import (
	"errors"
	"fmt"
	"genql/prismaUtil"
	"os"
	"strings"
)

func createCtx() {
	pathName := "./src/resolvers/context.ts"
	if !checkFileExists(pathName) {
		f, err := os.OpenFile(pathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		imports := "import { PrismaClient } from \"@prisma/client\";\nimport { Request, Response } from \"express\";"
		ctx := "export interface context {\n\tprisma: PrismaClient;\n\treq: Request;\n\tres: Response;\n}"
		f.WriteString(imports + "\n" + ctx)
	}
}
func getIdType(model *prismaUtil.Model) string {
	for _, f := range model.Fields {
		if f.Attribute != "" && strings.Index(f.Attribute, "@id") != -1 {
			typename, ok := prismaUtil.MAPPED_TS[f.Typename]
			if !ok {
				fmt.Println("id typename cannot be handled")
				os.Exit(1)
			}
			return typename
		}
	}
	panic("id typename cannot be handled")
}
func CreateResolver(model prismaUtil.Model) {
	addTypes(model)
	createCtx()

	objectType := model.Name
	createInputType := "create" + model.Name + "Input"
	updateInputType := "update" + model.Name + "Input"
	headers := "import { Arg, Ctx, Mutation, Query, Resolver } from \"type-graphql\";\n" +
		"import { context } from \"../context\"\n" +
		"import { " + objectType + ", " + createInputType + ", " + updateInputType + " } from \"./types\"\n\n"
	resolverClass := "@Resolver()\nexport class " + model.Name + "Resolver {"

	idType := getIdType(&model)
	getQuery := "\t@Query(() => " + objectType + ", { nullable: true })\n\tget" + model.Name + "(@Ctx() { prisma }: context, @Arg(\"id\") id: " + idType + "){\n"
	getFirstQuery := "\t\treturn prisma." + strings.ToLower(model.Name) + ".findFirst({\n\t\t\twhere: {\n\t\t\t\tid: id\n\t\t\t},\n\t\t});\n\t}"

	createMutation := "\t@Mutation(() => " + objectType + ")\n\tcreate" + model.Name + "(@Ctx() { prisma }: context, @Arg(\"input\") input: " + createInputType + "){\n"
	createQuery := "\t\treturn prisma." + strings.ToLower(model.Name) + ".create({\n\t\t\tdata: {\n\t\t\t\t...input\n\t\t\t},\n\t\t});\n\t}"

	updateMutation := "\t@Mutation(() => " + objectType + ")\n\tupdate" + model.Name + "(@Ctx() { prisma }: context, @Arg(\"input\") input: " + updateInputType + "){\n"
	updateQuery := "\t\treturn prisma." + strings.ToLower(model.Name) + ".update({\n\t\t\twhere:{id: input.id},\n\t\t\tdata: {\n\t\t\t\t...input\n\t\t\t},\n\t\t});\n\t}"

	deleteMutation := "\t@Mutation(() => " + objectType + ", { nullable: true })\n\tdelete" + model.Name + "(@Ctx() { prisma }: context, @Arg(\"id\") id: " + idType + "){\n"
	deleteQuery := "\t\treturn prisma." + strings.ToLower(model.Name) + ".delete({\n\t\t\twhere: {\n\t\t\t\tid: id\n\t\t\t},\n\t\t});\n\t}"

	resolverPath := "./src/resolvers/" + model.Name

	filePath := resolverPath + "/index.ts"
	if checkFileExists(filePath) {
		fmt.Printf("file (%s) already exists\n", filePath)
		os.Exit(1)
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(headers + resolverClass + "\n" + getQuery + getFirstQuery + "\n\n" + createMutation + createQuery + "\n\n" + updateMutation + updateQuery + "\n\n" + deleteMutation + deleteQuery + "\n\n" + "}")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addTypes(model prismaUtil.Model) {
	resolverPath := "./src/resolvers/" + model.Name

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
