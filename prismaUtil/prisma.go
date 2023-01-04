package prismaUtil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	case NPType:
		return "", nil
	}
	return "", errors.New("Invalid type")
}

type Field struct {
	Name       string
	Typename   PrismaType
	IsOptional bool
	IsArray    bool
	Attribute  string
	NPType     string // non-primative types, optional
}

type Model struct {
	Name   string
	Fields []Field
}

func (p *Model) String() string {
	stringVal := "\nmodel " + p.Name + " {\n"
	for _, field := range p.Fields {
		stringVal += "\t" + field.String() + "\n"
	}
	stringVal += "}"

	return stringVal
}

func (p *Model) AddField(field Field) {
	p.Fields = append(p.Fields, field)
}
func (p *Field) String() string {
	var prismaType string
	if p.Typename == NPType {
		prismaType = p.NPType
	} else {
		var err error
		prismaType, err = p.Typename.String()
		if err != nil {
			fmt.Println("invalid type entered")
			os.Exit(1)
		}
	}
	if p.IsArray {
		prismaType += "[]"
	}
	if p.IsOptional {
		prismaType += "?"
	}

	return p.Name + " " + prismaType + " " + p.Attribute
}

func parseID(values []string) PrismaType {
	if values[1] == "id" && len(values) == 3 {
		if values[2] == "uuid" {
			return StringType
		} else if values[2] == "ai" {
			return IntType
		}
		fmt.Printf("invalid default value for id type (%s)\n", values[2])
		os.Exit(1)
	}

	fmt.Printf("invalid id type entered (%s)\n", strings.Join(values, ":"))
	os.Exit(1)
	panic("")
}

func ParseField(str string) Field { // string of the form typename:type:default_value
	values := strings.Split(str, ":")
	if len(values) < 2 || len(values) > 3 {
		fmt.Printf("Invalid format enetered (%s)\n", strings.Join(values, ":"))
		os.Exit(1)
	}
	parsedT := Field{Name: values[0], IsOptional: false, IsArray: false, Attribute: ""}
	splitType := strings.Split(values[1], "[]")
	if len(splitType) == 2 {
		parsedT.IsArray = true
	}
	splitType = strings.Split(strings.Join(splitType, ""), "?")
	if len(splitType) == 2 {
		parsedT.IsOptional = true
	}

	if splitType[0][0] >= 65 && splitType[0][0] <= 90 {
		parsedT.NPType = splitType[0]
		parsedT.Typename = NPType
	} else {
		typename, ok := MAPPED_TYPES[splitType[0]]
		if !ok {
			// check if its an id type
			if values[1] == "id" {
				parsedT.Typename = parseID(values)
				parsedT.Attribute += "@id\t"
			} else {
				fmt.Printf("invalid type entered (%s), please enter a valid type\n", splitType[0])
				os.Exit(1)
			}
		} else {
			parsedT.Typename = typename
		}
	}

	// parse attributes
	if len(values) == 3 {
		if attribute, ok := MAPPED_ATTRIB[values[2]]; ok {
			parsedT.Attribute += attribute
		} else {
			parsedT.Attribute = "@default (\"" + values[2] + "\")"
		}
	}

	return parsedT
}

func ParseModel(modelName string, values []string) Model {
	exists := findModel(modelName)
	if exists {
		fmt.Printf("Model (%s) already exists\n", modelName)
		os.Exit(1)
	}
	parsedM := Model{Name: modelName, Fields: []Field{}}
	for _, val := range values {
		parsedM.Fields = append(parsedM.Fields, ParseField(val))
	}
	return parsedM
}

func GetIdType(modelName string) PrismaType {
	f, err := os.ReadFile(GetSchemaPath())
	if err != nil {
		fmt.Println("file not found. Create a prisma.schema file @ the following path: ", GetSchemaPath())
		os.Exit(1)
	}
	index := bytes.Index(f, []byte(modelName))

	index2 := bytes.Index(f[index:], []byte("@id"))
	var IdType string
	for i := index2 + index; i >= 0; i-- {
		if (f[i] >= 65 && f[i] <= 90) || (f[i] >= 97 && f[i] <= 122) {
			IdType = string(f[i]) + IdType
		} else if len(IdType) > 0 {
			break
		}
	}

	typename, ok := MAPPED_TYPES[strings.ToLower(IdType)]
	if !ok {
		fmt.Printf("unknown Id type (%s)\n", IdType)
		os.Exit(1)
	}
	return typename
}

func AddField(field Field, modelName string) {
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
		if (buffer[0] >= 65 && buffer[0] <= 90) || (buffer[0] >= 97 && buffer[0] <= 122) {
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
		fmt.Printf("Model name %s does not exist in schema.prisma, make sure to create the model first\n", modelName)
		os.Exit(1)
	}
	for err != io.EOF {
		i, err = f.Read(buffer)
		if buffer[0] == 125 {
			newBytes = append(newBytes, []byte("\t"+field.String()+"\n")...)
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

func findModel(modelName string) bool {
	f, err := os.ReadFile(GetSchemaPath())
	if err != nil {
		fmt.Printf("schema.prisma does not exist in %s", GetSchemaPath())
		os.Exit(1)
	}
	index := bytes.Index(f, []byte("model "+modelName))
	if index == -1 {
		return false
	}
	return true
}
func GetModel(modelName string) Model {
	f, err := os.ReadFile(GetSchemaPath())
	if err != nil {
		fmt.Printf("schema.prisma does not exist in %s", GetSchemaPath())
		os.Exit(1)
	}
	model := Model{Name: modelName, Fields: []Field{}}
	index := bytes.Index(f, []byte("model "+modelName))
	if index == -1 {
		fmt.Printf("Model (%s) not found in prisma.schema\n", modelName)
		os.Exit(1)
	}
	index2 := bytes.Index(f[index:], []byte("\n"))

	// var char byte
	element := ""
	for _, b := range f[index+index2:] {

		if b == '}' {
			break
		}
		if b != 10 && b != 32 {
			element += string(b)
		}
		if b == ' ' && len(element) > 0 && element[len(element)-1] != '-' {
			element += "-"
		}
		if b == 10 {
			if len(element) > 0 {
				model.Fields = append(model.Fields, parseField2(element))
			}
			element = ""

		}
	}

	return model
}

// parse field seperated by dashed value
func parseField2(dValue string) Field {
	values := strings.Split(dValue, "-")
	if len(values) < 2 {
		fmt.Println("error parsing values")
		os.Exit(1)
	}
	field := Field{Name: values[0]}
	splitType := strings.Split(values[1], "[]")
	if len(splitType) == 2 {
		field.IsArray = true
	}
	splitType = strings.Split(strings.Join(splitType, ""), "?")
	if len(splitType) == 2 {
		field.IsOptional = true
	}

	lowT := strings.ToLower(splitType[0])
	if t, ok := MAPPED_TYPES[lowT]; ok {
		field.Typename = t
	} else {
		field.Typename = NPType
		field.NPType = splitType[0]
	}

	if len(values) >= 3 {
		attrib := ""
		for _, att := range values[2:] {
			attrib += att + " "
		}
		field.Attribute = attrib
	}

	return field
}
