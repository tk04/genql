package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"genql/prismaUtil"
	"github.com/spf13/cobra"
)

// availible types in cli: date, int, string, json, bigint, bool, bytes, id, float
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate a Prisma Model",
	Long:  "Generate a Prisma model that is appended to the end of the schema.prisma file.\n\n Usage: model [model name] [list name:type:default_value].\n Example: genql model Test name:string id:id:ai isAdmin:bool:false",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oto, _ := cmd.Flags().GetString("OneToOne")
		otm, _ := cmd.Flags().GetString("OneToMany")
		mto, _ := cmd.Flags().GetString("ManyToOne")
		fmt.Println("oto: ", oto, " -- otm: ", otm, "-- mto: ", mto)

		prismaModel := prismaUtil.ParseModel(args[0], args[1:])
		f, err := os.OpenFile(GetSchemaPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(prismaModel.String())
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
		"ai":     "@default(autoincrement())",
		"uuid":   "String\t\t @id @default(uuid())",
		"true":   "@default(true)",
		"false":  "@default(false)",
		"unique": "@unique",
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

func BuildOneToOne() {
	addField("testField\tString", "User")
}

func getIdType(modelName string) string {

	f, err := os.ReadFile(GetSchemaPath())
	if err != nil {
		fmt.Println("file not found. Create a prisma.schema file @ the following path: ", GetSchemaPath())
		os.Exit(1)
	}
	index := bytes.Index(f, []byte(modelName))
	fmt.Println(index)

	index2 := bytes.Index(f[index:], []byte("@id"))
	fmt.Println(index2)
	var IdType string
	for i := index2 + index; i >= 0; i-- {
		if isUpper(f[i]) || (f[i] >= 97 && f[i] <= 122) {
			IdType = string(f[i]) + IdType
		} else if len(IdType) > 0 {
			break
		}
	}
	return IdType
}

func addField(field string, modelName string) {
	var newBytes []byte

	f, err := os.OpenFile(GetSchemaPath(), os.O_RDWR, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("file not found. Create a schema.prisma file @ the following path: ", GetSchemaPath())
		os.Exit(1)
	}
	buffer := make([]byte, 1)
	i, err := f.Read(buffer)
	offset := i
	currStr := ""
	for err != io.EOF {
		newBytes = append(newBytes, buffer...)
		if isUpper(buffer[0]) || (buffer[0] >= 97 && buffer[0] <= 122) {
			currStr += string(buffer)
		} else {
			if currStr == modelName {
				break
			}
			currStr = ""
		}
		i, err = f.Read(buffer)
		offset += i
	}
	if err == io.EOF {
		fmt.Printf("Model name %s does not exist in schema.prisma, make sure to create the model first", modelName)
		os.Exit(1)
	}
	for err != io.EOF {
		i, err = f.Read(buffer)
		if buffer[0] == 125 {
			newBytes = append(newBytes, []byte("\t"+field+"\n")...)
			break
		}
		newBytes = append(newBytes, buffer...)
	}
	for err != io.EOF {
		newBytes = append(newBytes, buffer...)
		i, err = f.Read(buffer)
	}
	f.Truncate(0)
	f.Seek(0, 0)
	f.Write(newBytes)
}
