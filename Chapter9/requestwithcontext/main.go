package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ClientContext struct {
	http.Client
}

func (c *ClientContext) GetContext(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func main() {
	client := &ClientContext{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.GetContext(ctx, "http://www.example.com")
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	fmt.Println(string(bytes))
}
