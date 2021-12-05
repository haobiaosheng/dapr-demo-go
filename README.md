# 介绍
这是一个简单的dapr golang http调用demo，包括一个service端和client端。

# 版本介绍
- go版本：1.14.6
- dapr go sdk版本：v1.3.0

# 工程目录
```bash
.
├── client
│   └── client.go
├── go.mod
├── go.sum
├── internal
│   └── response.go
└── service
    └── service.go
```

- service包：服务端代码
- client包：客户端代码
- internal包：内部包

# 服务端代码
> 比较简单，以9003端口启动app，并监听/hello url。

```go
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
```

# 服务端启动
```bash
dapr run --app-id go-service \
         --app-protocol http \
         --app-port 9003 \
         --dapr-http-port 3501 \
         --log-level debug \
         go run ./service/service.go
```

- app-id：指定用于服务发现的应用程序id
- app-protocol：dapr用于与应用程序通信，有效值为: http或grpc。
- app-port：应用程序正在侦听的端口，这里指定端口为9003。
- dapr-http-port：dapr本身要监听的HTTP端口，用于给外部调用。
- log-lever：日志级别，有效值因为其中之一: debug, info, warn, error, fatal, or panic。

启动成功后日志如下：
```bash
ℹ️  Updating metadata for app command: go run ./service/service.go
✅  You're up and running! Both Dapr and your app logs will appear here.
```

# 客户端代码
> 启动一个client，绑定服务端app-id和method，进行通信调用。

```go
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
```

# 客户端启动
```bash
dapr run --app-id go-client \
         --log-level debug \
         go run ./client/client.go
```

启动成功后的日志如下：
```bash
ℹ️  Dapr sidecar is up and running.
ℹ️  Updating metadata for app command: go run ./client/client.go
✅  You're up and running! Both Dapr and your app logs will appear here.
```

client端由于设置了启动每隔5s就去访问server端，可以看到两边的日志都有输出。

client端的日志输出
```bash
== APP == 2021/12/05 17:34:42 go-service method hello has invoked, response is: {"Message":"This message is from service"}
DEBU[0004] mDNS browse for app id go-service timed out.  app_id=go-client instance=hbshong-dev scope=dapr.contrib type=log ver=1.5.0
DEBU[0004] Refreshing mDNS addresses for app id go-service timed out.  app_id=go-client instance=hbshong-dev scope=dapr.contrib type=log ver=1.5.0
DEBU[0006] found mDNS IPv4 address in cache: 9.134.6.54:36864  app_id=go-client instance=hbshong-dev scope=dapr.contrib type=log ver=1.5.0
== APP == 2021/12/05 17:34:47 go-service method hello has invoked, response is: {"Message":"This message is from service"}
```

server端的日志输出
```bash
== APP == 2021/12/05 17:34:57 the service method hello has invoked, receive message is This is client
== APP == 2021/12/05 17:35:02 the service method hello has invoked, receive message is This is client
```

从两边的日志可以看到，client和server通过dapr实现了交互。