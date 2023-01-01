package cmd

import (
	"fmt"
	"genql/prismaUtil"
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
		mto, _ := cmd.Flags().GetString("ManyToOne")
		fmt.Println("oto: ", oto, " -- otm: ", otm, "-- mto: ", mto)

		prismaModel := prismaUtil.ParseModel(args[0], args[1:])
		f, err := os.OpenFile(prismaUtil.GetSchemaPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(prismaModel.String())
		if err != nil {
			fmt.Println(err)
		}
	},
}

func buildOneToOne(fModelName string, sModelName string) {
	field := prismaUtil.Field{Name: strings.ToLower(fModelName), IsOptional: true, IsArray: false, Typename: prismaUtil.NPType, NPType: fModelName}
	prismaUtil.AddField(field, sModelName)
}
