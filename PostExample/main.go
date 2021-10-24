package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const json = `{"name": "Matt", "age": 45}`		// This is our JSON

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	in := strings.NewReader(json)				// Wrap JSON with an io.Reader

	// Issue HTTP POST, declaring our content-type as "text/json"
	resp, err := client.Post("https://localhost:8080/v1/c8/key-e","text/json", in)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(message))
}
