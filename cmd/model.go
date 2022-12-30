package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// availible types in cli: date, int, string, json, bigint, bool, bytes, id, float
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate a Prisma Model",
	Long:  "Generate a Prisma model that is appended to the end of the schema.prisma file.\n\n Usage: model [model name] [list name:type:default_value].\n Example: genql model Test name:string id:id:ai isAdmin:bool:false",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		modelName := args[0]
		mappedTypes := mapType(args[1:])
		fmt.Println(mappedTypes)
		modelString := "\nmodel " + modelName + " {\n"
		for k, v := range mappedTypes {
			modelString += "\t" + k + "\t\t\t\t" + v + "\n"
		}
		modelString += "}"
		f, err := os.OpenFile(GetSchemaPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(modelString)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func GetSchemaPath() string {
	cmd := exec.Command("pwd")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	path := out.String()
	return path[:len(path)-1] + "/prisma/schema.prisma"
}

func isUpper(char byte) bool {
	return char >= 65 && char <= 90
}
func mapType(values []string) map[string]string {
	mappedAtrrib := map[string]string{
		"ai":    "@default(autoincrement())",
		"uuid":  "String\t\t @id @default(uuid())",
		"true":  "@default(true)",
		"false": "@default(false)",
	}

	mappedTypes := map[string]string{
		"string": "String",
		"int":    "Int",
		"bigint": "BigInt",
		"float":  "Float",
		"id":     "Int\t\t @id",
		"date":   "DateTime",
		"json":   "Json",
		"bytes":  "Bytes",
		"bool":   "Boolean",
	}
	mappedValues := make(map[string]string)
	// if value name starts w/ an uppercase char, put the value as it is
	for _, str := range values {
		values := strings.Split(str, ":")
		fmt.Println(values)
		//handle attributes
		var attribute string
		if len(values) == 3 {
			if _, ok := mappedAtrrib[values[2]]; ok {
				attribute = mappedAtrrib[values[2]]
				if values[2] == "uuid" {
					mappedValues[values[0]] = attribute
					continue
				}
			} else {
				attribute = "@default (\"" + values[2] + "\")"
			}

		}
		// handle fields
		if isUpper(values[1][0]) {
			mappedValues[values[0]] = values[1]
			continue
		}
		_, ok := mappedTypes[values[1]]
		if !ok {
			fmt.Printf("Invalid type: (%s)\n", values[1])
			os.Exit(1)
		}
		mappedValues[values[0]] = mappedTypes[values[1]] + " " + attribute
	}
	return mappedValues
}
