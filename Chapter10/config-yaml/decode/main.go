package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type TaggedMore struct {
	// Flow will be marshalled using a "flow" style
	// (useful for structs, sequences and maps).
	Flow map[string]string `yaml:"flow"`

	// Inlines a struct or a map, causing all of its fields
	// or keys to be processed as if they were part of the outer
	// struct. For maps, keys must not conflict with the yaml
	// keys of other struct fields.
	Inline map[string]string `yaml:",inline"`
}

type Config struct {
	Host string
	Port uint16
	Tags map[string]string
}

func main() {
	c := Config{}
	//	bytes := []byte(`
	//host: 127.0.0.1
	//port: 1234
	//tags:
	//    foo: bar
	//`)

	file, err := os.Open("CloudNativeGo/Chapter10/config-yaml/decode/sample.yaml")
	defer file.Close()
	body, err := io.ReadAll(file)
	err = yaml.Unmarshal(body, &c)
	//err := yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Printf("couldn't unmarshal: %w", err)
	}
	fmt.Println(c)
}
