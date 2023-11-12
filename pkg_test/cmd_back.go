package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cmdback = &cobra.Command{
	Use:   "backup  [flags] [FILE/DIR] ...",
	Short: "test short backup",
	Long:  "test Long backup",
	PreRun: func(cmd *cobra.Command, args []string) {
		hostname, err := os.Hostname()
		if err != nil {
			return
		}
		fmt.Println(hostname)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("backup", cmd, args)
	},
}

func init() {
	fmt.Println(123)
	cmdroot.AddCommand(cmdback)
	//f := cmdback.Flags()
	//par := "root"
	//f.StringVar(&par, "S", "s", "s")
}
