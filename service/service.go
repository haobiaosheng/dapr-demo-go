package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"dapr-golang/internal"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

func main() {
    s := daprd.NewService(":9003")
    if err := s.AddServiceInvocationHandler("/hello", helloHandler); err != nil {
        log.Fatalf("add invocation handler err, the err is: %v", err)
    }
    if err := s.Start(); err != nil  && err != http.ErrServerClosed {
        log.Fatalf("listening err, the err is: %v", err)
    }
}

func helloHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
    if in == nil {
        return nil, errors.New("invocation parameter required")
    }

    log.Printf("the whole input is: %+v\n", in)
    log.Printf("the service method hello has invoked, receive message is %v\n", string(in.Data))
    resp := internal.HTTPResp{
        Message: "This message is from service",
    }
    return &common.Content{
		Data:        resp.ToBytes(),
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}, nil
}