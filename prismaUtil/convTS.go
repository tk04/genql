package prismaUtil

import (
	"fmt"
	"os"
	"strings"
)

var MAPPED_TS = map[PrismaType]string{
	FloatType:    "number",
	IntType:      "number",
	BytesType:    "Uint8Array",
	JsonType:     "any",
	BooleanType:  "boolean",
	StringType:   "string",
	BigIntType:   "number",
	DateTimeType: "Date",
}

func (m Model) toTS(optional bool) string {
	tsType := ""
	for _, field := range m.Fields {
		// skip non-primative types
		if field.Typename == NPType {
			continue
		}

		if optional && strings.Index(field.Attribute, "@id") != -1 { // id required for update operation
			tsType += "\t@Field()\n\t" + field.Name
		} else if optional || field.IsOptional || strings.Index(field.Attribute, "@default") != -1 {
			tsType += "\t@Field({ nullable: true })\n\t" + field.Name
			tsType += "?"
		} else {
			tsType += "\t@Field()\n\t" + field.Name
		}
		tsType += ": "

		typename, ok := MAPPED_TS[field.Typename]
		if !ok {
			fmt.Println("invalid type encountered")
			os.Exit(1)
		}
		tsType += typename

		if field.IsArray {
			tsType += "[]"
		}
		tsType += "\n"
	}
	tsType += "}"
	return tsType
}
func (m Model) ObjectType() string {
	objectType := "@ObjectType()\nexport class " + m.Name + "{\n"
	objectType += m.toTS(false)

	return objectType
}
func (m Model) CreateInputType() string {
	inputType := "@InputType()\nexport class " + "create" + m.Name + "Input " + "{\n"
	inputType += m.toTS(false)

	return inputType
}

func (m Model) UpdateInputType() string {
	inputType := "@InputType()\nexport class " + "update" + m.Name + "Input " + "{\n"

	inputType += m.toTS(true)
	return inputType
}
