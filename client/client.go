package main

import (
	"context"
	"log"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

func main() {
    ctx := context.Background()

    // create the dapr
    client, err := dapr.NewClient()
    if err != nil {
        panic(err)
    }
    defer client.Close()

    content := &dapr.DataContent{
		Data:        []byte("This is client"),
		ContentType: "text/plain",
	}
    for {
        resp, err := client.InvokeMethodWithContent(ctx, "go-service", "hello", "get", content)
        if err != nil {
            panic(err)
        }
        log.Printf("go-service method hello has invoked, response is: %s", string(resp))
        time.Sleep(time.Second * 5)
    }
}