package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type TypeToken uint8

const (
	StringToken TypeToken = iota
	IntToken
	DateTimeToken
	IdToken
	BooleanType
	FloatType
	BigIntType
	JsonType
	BytesType
)

func (TypeToken) Error(valType string) string {
	return fmt.Sprintf("Invalid Type: %s", valType)
}
func (p TypeToken) String() (string, error) {
	switch p {
	case StringToken:
		return "String", nil
	case IntToken:
		return "Int", nil
	case DateTimeToken:
		return "DateTime", nil
	case IdToken:
		return "ID", nil
	case BooleanType:
		return "Boolean", nil
	case FloatType:
		return "Float", nil
	case BigIntType:
		return "BigInt", nil
	case JsonType:
		return "Json", nil
	case BytesType:
		return "Bytes", nil
	}
	return "", errors.New("Invalid type")
}

// availible types in cli: date, int, string, json, bigint, bool, bytes, id, float
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate Prisma Model",
	Long:  "Generate a Prisma model that is appended to the end of the schema.prisma file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// modelName := args[0]
		// fmt.Println(modelName)
		// fmt.Println(mappedValues)
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

func mapType(values []string) map[string]TypeToken {

	mappedTypes := map[string]TypeToken{
		"string": StringToken,
		"int":    IntToken,
		"bigint": BigIntType,
		"float":  FloatType,
		"id":     IdToken,
		"date":   DateTimeToken,
		"json":   JsonType,
		"bytes":  BytesType,
		"bool":   BooleanType,
	}
	// schemaPath := GetSchemaPath()
	mappedValues := make(map[string]TypeToken)
	// if name starts w/ uppercase char, put the value as it is ?
	for _, str := range values {
		values := strings.Split(str, ":")
		_, ok := mappedTypes[values[1]]
		if !ok {
			fmt.Printf("Invalid type: (%s)\n", values[1])
			os.Exit(1)
		}
		mappedValues[values[0]] = mappedTypes[values[1]]
	}
	return mappedValues
}
