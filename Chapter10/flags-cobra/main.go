package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var strp string
var intp int
var boolp bool


var rootCmd = &cobra.Command{
	Use:  "cng",
	Long: "A super simple command.",
}

var flagsCmd = &cobra.Command{
	Use: "flags",
	Long: "A simple flags experimentation command, built with Cobra.",
	Run: flagsFunc,
}

func init() {
	flagsCmd.Flags().StringVarP(&strp, "string", "s", "foo", "a string")
	flagsCmd.Flags().IntVarP(&intp, "number", "n", 42, "an integer")
	flagsCmd.Flags().BoolVarP(&boolp,"boolean","b",false,"a boolean")

	rootCmd.AddCommand(flagsCmd)
}


func flagsFunc(cmd *cobra.Command, args []string) {
	fmt.Println("string:", strp)
	fmt.Println("integer:", intp)
	fmt.Println("boolean:", boolp)
	fmt.Println("args:", args)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}