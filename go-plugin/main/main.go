package main

import (
	"fmt"
	"log"
	"os"
	"plugin"
)

type Sayer interface {
	Says() string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: run main/main.go animal")
	}
	// Get the animal name, and build the path where we expect to
	// find the corresponding shared object (.so) file.
	name := os.Args[1]
	module := fmt.Sprintf("./%s/%s.so", name, name)
	// Open our plugin and get a *plugin.Plugin.
	p, err := plugin.Open(module)
	if err != nil {
		log.Fatal(err)
	}
	// Lookup searches for a symbol named "Animal" in plug-in p.
	symbol, err := p.Lookup("Animal")
	if err != nil {
		log.Fatal(err)
	}
	// Asserts that the symbol interface holds a Sayer.
	animal, ok := symbol.(Sayer)
	if !ok {
		log.Fatal("that's not a Sayer")
	}

	// Now we can use our loaded plug-in!
	fmt.Printf("A %s says: %q\n", name, animal.Says())
}
