package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	Host string
	Port uint16
	Tags map[string]string
}

func main() {
	c := Config{}
	//bytes := []byte(`{"Host":"127.0.0.1","Port":1234,"Tags":{"foo":"bar"}}`)
	file, err := os.Open("CloudNativeGo/Chapter10/config-json/decode/sample.json")
	defer file.Close()
	body, err := io.ReadAll(file)
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Printf("couldn't unmarshal: %w", err)
	}
	fmt.Println(c)
}
