package prismaUtil

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var MAPPED_ATTRIB = map[string]string{
	"ai":     "@default(autoincrement())",
	"uuid":   "@default(uuid())",
	"true":   "@default(true)",
	"false":  "@default(false)",
	"unique": "@unique",
}

var MAPPED_TYPES = map[string]PrismaType{
	"string": StringType,
	"int":    IntType,
	"bigint": BigIntType,
	"float":  FloatType,
	"date":   DateTimeType,
	"json":   JsonType,
	"bytes":  BytesType,
	"bool":   BooleanType,
}

type PrismaType uint8

const (
	StringType PrismaType = iota
	IntType
	DateTimeType
	BooleanType
	FloatType
	BigIntType
	JsonType
	BytesType
	NPType // non-primative types
)

func (p PrismaType) String() (string, error) {
	switch p {
	case StringType:
		return "String", nil
	case IntType:
		return "Int", nil
	case DateTimeType:
		return "DateTime", nil
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

type Field struct {
	name       string
	typename   PrismaType
	isOptional bool
	isArray    bool
	attribute  string
	npType     string // non-primative types, optional
}

type Model struct {
	name   string
	fields []Field
}

func (p *Model) String() string {
	stringVal := "\nmodel " + p.name + " {\n"
	for _, field := range p.fields {
		stringVal += "\t" + field.String() + "\n"
	}
	stringVal += "}"

	return stringVal
}

func (p *Field) String() string {
	prismaType, err := p.typename.String()
	if err != nil {
		fmt.Println("invalid type entered")
		os.Exit(1)
	}
	if p.isArray {
		prismaType += "[]"
	}
	if p.isOptional {
		prismaType += "?"
	}

	return p.name + "\t\t" + prismaType + "\t" + p.attribute
}

func parseID(values []string) PrismaType {
	if values[1] == "id" {
		if values[2] == "uuid" {
			return StringType
		} else if values[2] == "ai" {
			return IntType
		}
		fmt.Printf("invalid default value for id type (%s)\n", values[2])
		os.Exit(1)
	}

	panic("invalid id type entered")
}

func ParseField(str string) Field { // string of the form typename:type:default_value
	values := strings.Split(str, ":")
	parsedT := Field{name: values[0], isOptional: false, isArray: false, attribute: ""}
	splitType := strings.Split(values[1], "[]")
	if len(splitType) == 2 {
		parsedT.isArray = true
	}
	splitType = strings.Split(strings.Join(splitType, ""), "?")
	if len(splitType) == 2 {
		parsedT.isOptional = true
	}

	if splitType[0][0] >= 65 || splitType[0][0] <= 90 {
		parsedT.npType = splitType[0]
		parsedT.typename = NPType
	}
	typename, ok := MAPPED_TYPES[splitType[0]]
	if !ok {
		// check if its an id type
		if values[1] == "id" {
			parsedT.typename = parseID(values)
			parsedT.attribute += "@id\t"
		} else {
			fmt.Printf("invalid type entered (%s), please enter a valid type\n", splitType[0])
			os.Exit(1)
		}
	} else {
		parsedT.typename = typename
	}

	// parse attributes
	if len(values) == 3 {
		if attribute, ok := MAPPED_ATTRIB[values[2]]; ok {
			parsedT.attribute += attribute
		} else {
			parsedT.attribute = "@default (\"" + values[2] + "\")"
		}
	}

	return parsedT
}

func ParseModel(modelName string, values []string) Model {
	parsedM := Model{name: modelName, fields: []Field{}}
	for _, val := range values {
		parsedM.fields = append(parsedM.fields, ParseField(val))
	}
	return parsedM
}
