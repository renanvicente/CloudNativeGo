package main

import (
	"flag"
	"fmt"
)

func main() {
	// Declare a string flag with a default value "foo"
	// and a short description. It returns a string pointer.
	strp := flag.String("string","foo","a string")

	// Declare number and Boolean flags, similar to the string flag.
	intp := flag.Int("number",42,"an integer")
	boolp := flag.Bool("boolean",false, "a boolean")

	// Call flag.Parse() to execute command-line parsing.
	flag.Parse()

	// Print the parsed options and trailing positional arguments.
	fmt.Println("string:", *strp)
	fmt.Println("integer:", *intp)
	fmt.Println("boolean:", *boolp)
	fmt.Println("args:", flag.Args())

}
