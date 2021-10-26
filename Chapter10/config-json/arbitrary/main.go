package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	bytes := []byte(`{"Foo":"Bar", "Number":1313, "Tags":{"A":"B"}}`)
	var f interface{}
	err := json.Unmarshal(bytes, &f)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(f)
}
