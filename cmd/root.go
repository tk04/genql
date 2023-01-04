package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genql",
	Short: "genql is a server side GraphQL & Prisma code generator",
	Long:  "A code generator that reliably generates Prisma database schemas followed by CRUD GraphQL resolvers for each generated model.",
}

func Execute() {
	rootCmd.AddCommand(modelCmd)
	rootCmd.AddCommand(resolversCmd)

	var OTMRelation string // one to many relationship
	var OTORelation string // one to one relationship
	var MTORelation string // many to one relationship
	modelCmd.Flags().StringVarP(&OTMRelation, "OneToMany", "r", "", "Define a one-to-many relationship between two models")
	modelCmd.Flags().StringVarP(&OTORelation, "OneToOne", "1", "", "Define a one-to-one relationship between two models")
	modelCmd.Flags().StringVarP(&MTORelation, "ManyToMany", "m", "", "Define a many-to-one relationship between two models")

	var Exceptions []string
	resolversCmd.Flags().StringArrayVarP(&Exceptions, "Except", "e", []string{}, "Define operations not to be included in a given resolver")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
