package main

import (
	"fmt"
	"github.com/spf13/cobra"
)
import "context"

var cmdroot = &cobra.Command{
	Use:   "root [flags]",
	Short: "test short",
	Long:  "test Long",

	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		fmt.Println("rootper", c, args)
		return nil
	},
}
var internalGlobalCtx context.Context

func main() {
	kt := Kt
	fmt.Println(kt)
	fmt.Println(234)
	err := cmdroot.Execute()
	if err != nil {
		fmt.Errorf("err")
	}

}
