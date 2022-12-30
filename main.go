package main

import (
	"fmt"
	"genql/cmd"
)

func main() {
	cmd.Execute()
	fmt.Println(cmd.GetSchemaPath())
}
