package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Tagged struct {
	// CustomKey will appear in JSON as the key "custom_key".
	CustomKey string `json:"custom_key"`

	// OmitEmpty will appear in JSON as "OmitEmpty" (the default),
	// but will only be written if it contains a nonzero value.
	OmitEmpty string `json:",omitempty"`

	// IgnoredName will always be ignored.
	IgnoredName string `json:"-"`

	// TwoThings will appear in JSON as the key "two_things",
	// but only if it isn't empty.
	TwoThings string `json:"two_things,omitempty"`
}

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

	//bytes, err := json.Marshal(c)
	bytes, err := json.MarshalIndent(c, "", "    ")

	if err != nil {
		log.Printf("couldn't marshal: %w", err)
	}
	fmt.Println(string(bytes))
}
