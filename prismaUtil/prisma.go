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

	return p.Name + "\t\t" + prismaType + "\t" + p.Attribute
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
	parsedM := Model{name: modelName, fields: []Field{}}
	for _, val := range values {
		parsedM.fields = append(parsedM.fields, ParseField(val))
	}
	return parsedM
}

func GetIdType(modelName string) string {

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
		if (f[i] >= 65 && f[i] <= 90) || (f[i] >= 97 && f[i] <= 122) {
			IdType = string(f[i]) + IdType
		} else if len(IdType) > 0 {
			break
		}
	}
	return IdType
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
		fmt.Printf("Model name %s does not exist in schema.prisma, make sure to create the model first", modelName)
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
