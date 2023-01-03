package resolvers

import (
	"errors"
	"fmt"
	"genql/prismaUtil"
	"os"
	"strings"
)

type Resolver struct {
	Functions []string
	Model     prismaUtil.Model
}

func (r Resolver) CreateFiles() {
	addTypes(r.Model)
	createCtx()

	resolverPath := "./src/resolvers/" + r.Model.Name

	filePath := resolverPath + "/index.ts"
	if checkFileExists(filePath) {
		fmt.Printf("file (%s) already exists\n", filePath)
		os.Exit(1)
	}
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(r.String())
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
func (r Resolver) String() string {
	createInputType := "create" + r.Model.Name + "Input"
	updateInputType := "update" + r.Model.Name + "Input"
	headers := "import { Arg, Ctx, Mutation, Query, Resolver } from \"type-graphql\";\n" +
		"import { context } from \"../context\"\n" +
		"import { " + r.Model.Name + ", " + createInputType + ", " + updateInputType + " } from \"./types\"\n\n"
	resolverClass := "@Resolver()\nexport class " + r.Model.Name + "Resolver {\n"
	ts := headers + resolverClass
	idType := getIdType(&r.Model)
	for _, val := range r.Functions {
		switch val {
		case "get":
			ts += addFunc(r.Model.Name, idType) + "\n"
		case "create":
			ts += createFunc(r.Model.Name) + "\n"
		case "update":
			ts += updateFunc(r.Model.Name) + "\n"
		case "delete":
			ts += deleteFunc(r.Model.Name, idType) + "\n"
		}
	}

	ts += "}"
	return ts
}
func addFunc(modelName string, idType string) string {
	// objectType := model.Name

	getQuery := "\t@Query(() => " + modelName + ", { nullable: true })\n\tget" + modelName + "(@Ctx() { prisma }: context, @Arg(\"id\") id: " + idType + "){\n"
	getFirstQuery := "\t\treturn prisma." + strings.ToLower(modelName) + ".findFirst({\n\t\t\twhere: {\n\t\t\t\tid: id\n\t\t\t},\n\t\t});\n\t}"
	return getQuery + getFirstQuery
}

func createFunc(modelName string) string {
	createInputType := "create" + modelName + "Input"
	createMutation := "\t@Mutation(() => " + modelName + ")\n\tcreate" + modelName + "(@Ctx() { prisma }: context, @Arg(\"input\") input: " + createInputType + "){\n"
	createQuery := "\t\treturn prisma." + strings.ToLower(modelName) + ".create({\n\t\t\tdata: {\n\t\t\t\t...input\n\t\t\t},\n\t\t});\n\t}"
	return createMutation + createQuery
}

func updateFunc(modelName string) string {
	updateInputType := "update" + modelName + "Input"

	updateMutation := "\t@Mutation(() => " + modelName + ")\n\tupdate" + modelName + "(@Ctx() { prisma }: context, @Arg(\"input\") input: " + updateInputType + "){\n"
	updateQuery := "\t\treturn prisma." + strings.ToLower(modelName) + ".update({\n\t\t\twhere:{id: input.id},\n\t\t\tdata: {\n\t\t\t\t...input\n\t\t\t},\n\t\t});\n\t}"
	return updateMutation + updateQuery
}

func deleteFunc(modelName string, idType string) string {
	deleteMutation := "\t@Mutation(() => " + modelName + ", { nullable: true })\n\tdelete" + modelName + "(@Ctx() { prisma }: context, @Arg(\"id\") id: " + idType + "){\n"
	deleteQuery := "\t\treturn prisma." + strings.ToLower(modelName) + ".delete({\n\t\t\twhere: {\n\t\t\t\tid: id\n\t\t\t},\n\t\t});\n\t}"

	return deleteMutation + deleteQuery
}

func FornatTS(lines []string) string {
	ts := ""
	val := struct{}{}
	seps := map[byte]struct{}{
		'{': val,
		'(': val,
		'[': val,
	}
	ends := map[byte]struct{}{
		'}': val,
		')': val,
		']': val,
	}

	tabs := []string{}
	for _, line := range lines {
		ws := strings.Trim(line, " ")
		ws = strings.Trim(ws, ",")
		ws = strings.Trim(ws, ";")
		if _, ok := ends[ws[len(ws)-1]]; ok && len(tabs) > 0 {
			tabs = tabs[: len(tabs)-1 : len(tabs)-1]
		}
		ts += strings.Join(tabs, "") + line + "\n"
		if _, ok := seps[ws[len(ws)-1]]; ok {
			tabs = append(tabs, "\t")
		}
	}
	return ts
}

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
