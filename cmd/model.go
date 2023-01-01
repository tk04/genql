package cmd

import (
	"fmt"
	"genql/prismaUtil"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"os"
	"strings"
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
		mtm, _ := cmd.Flags().GetString("ManyToMany")

		relations := []prismaUtil.Field{}
		if oto != "" {
			relations = append(relations, buildOneToOne(oto, args[0])...)
		}
		if otm != "" {
			relations = append(relations, buildOneToMany(otm, args[0])...)
		}
		if mtm != "" {
			relations = append(relations, buildManyToMany(mtm, args[0]))
		}
		prismaModel := prismaUtil.ParseModel(args[0], args[1:])

		for _, rel := range relations {
			prismaModel.AddField(rel)
		}

		f, err := os.OpenFile(prismaUtil.GetSchemaPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(prismaModel.String())
		if err != nil {
			fmt.Println(err)
		}
	},
}

func buildOneToOne(values string, fModelName string) []prismaUtil.Field {
	vals := strings.Split(values, ":")
	if len(vals) != 2 {
		fmt.Printf("Invalid relation format (%s)\n", values)
	}
	relationField := prismaUtil.Field{Name: vals[0], Attribute: "@relation(fields: [" + strings.ToLower(vals[1]) + "Id" + "], references: [id])", Typename: prismaUtil.NPType, NPType: vals[1]}
	idField := prismaUtil.Field{Name: strings.ToLower(vals[1]) + "Id", Typename: prismaUtil.GetIdType(vals[1]), Attribute: "@unique"}

	field := prismaUtil.Field{Name: strings.ToLower(fModelName), IsOptional: true, IsArray: false, Typename: prismaUtil.NPType, NPType: fModelName}
	prismaUtil.AddField(field, vals[1])
	return []prismaUtil.Field{relationField, idField}
}

func buildOneToMany(values string, fModelName string) []prismaUtil.Field {
	vals := strings.Split(values, ":")
	if len(vals) != 2 {
		fmt.Printf("Invalid relation format (%s)\n", values)
	}
	relationField := prismaUtil.Field{Name: vals[0], Attribute: "@relation(fields: [" + strings.ToLower(vals[1]) + "Id" + "], references: [id])", Typename: prismaUtil.NPType, NPType: vals[1]}
	idField := prismaUtil.Field{Name: strings.ToLower(vals[1]) + "Id", Typename: prismaUtil.GetIdType(vals[1])}

	field := prismaUtil.Field{Name: strings.ToLower(fModelName), IsArray: true, Typename: prismaUtil.NPType, NPType: fModelName}
	prismaUtil.AddField(field, vals[1])
	return []prismaUtil.Field{relationField, idField}
}

func buildManyToMany(values string, fModelName string) prismaUtil.Field {
	pluralize := pluralize.NewClient()

	vals := strings.Split(values, ":")
	if len(vals) != 2 {
		fmt.Printf("Invalid relation format (%s)\n", values)
	}
	relationField := prismaUtil.Field{Name: pluralize.Plural(strings.ToLower(vals[0])), IsArray: true, Typename: prismaUtil.NPType, NPType: vals[1]}

	field := prismaUtil.Field{Name: pluralize.Plural(strings.ToLower(fModelName)), IsArray: true, Typename: prismaUtil.NPType, NPType: fModelName}
	prismaUtil.AddField(field, vals[1])
	return relationField
}
