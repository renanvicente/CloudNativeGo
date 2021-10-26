package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
)

type Config struct {
	Host string
	Port uint16
	Tags map[string]string
}

func main() {
	c := Config{
		Host: "localhost",
		Port: 1313,
		Tags: map[string]string{"env": "dev"},
	}
	bytes, err := yaml.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(bytes))
}
