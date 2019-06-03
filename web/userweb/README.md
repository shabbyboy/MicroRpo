### 长连接和短连接编写方式

* 短连接


短连接直接沿用了micro web的做法，没有任何特殊处理，例如：

```$xslt
service := web.NewService(
    web.Name("microrpo.web.user"),
)
if err := service.Init(); err != nil{
    log.Fatal(err)
}
service.HandleFunc("/user/login", func(writer http.ResponseWriter, request *http.Request) {
    if request.Method == "GET"{
        request.ParseForm()
        name := request.Form.Get("name")
        if len(name) == 0{
            name = "world"
        }

        cl := msgproto.NewSayHelloService("microrpo.srv.websrv",client.DefaultClient)

        rsp, err := cl.Hello(context.Background(),&msgproto.Request{
            Name:name,
        })

        if err != nil{
            http.Error(writer,err.Error(),500)
            return
        }
        log.Print(request.Host)
        writer.Write([]byte(`<html><body><h1>` + rsp.Msg + `</h1></body></html>`))
        return
    }

    fmt.Fprint(writer,`<html><body><h1>Enter Name<h1><form method=post><input name=name type=text /></form></body></html>`)
})

//service.Options().Service.Run()

if err := service.Run(); err != nil{
    log.Fatal(err)
}
```

>  go run ./main.go

>  micro api --api_namespace=microrpo.web --handler=web

curl "http:ip:8080/user/login?name=zhangsan"

* 长连接

    * 回调方法，实际处理收到的tcp业务

    ```$xslt
    type pub struct {
    
    
    }
    //event.id 必须是userid
    func (p *pub) Process(ctx context.Context,event *pubsub.Event) error{
    
        //// 通过发布订阅，往长连接发送消息
        pubproj := pubsubproj.Publish{
            Ctx:context.Background(),
            Client:client.DefaultClient,
        }
        fmt.Println(event.Id)
        publisher := pubproj.NewPublisher(event.Id)
    
        if err := pubproj.PubEvent(publisher,event); err != nil{
            fmt.Println("发布：",err)
        }
        return nil
    }
    ```

    * 注册tcp 服务
    
    ```$xslt
    subproj := pubsubproj.Subscribe{
        Server:server.DefaultServer,
    }
    err := subproj.SubTopic("tcp",new(pub))
    //发现主动web服务的时候，订阅服务一直存在，所以主动取消下订阅的服务
    defer subproj.UnSubTopic()
    if err != nil{
        log.Fatal(err)
    }
    subproj.Run()
    ```
    * 启动服务：
    > go run ../tcpweb/main.go
    
    > go run ./main.go

    * 客户端调用：
        
        1. 发起连接请求(推荐用websocket在线测试工具，http://coolaf.com/tool/chattest)
           > ws://ip:8080/wstcp?userid=123 //userid 需要是能唯一标识客户端
      
        2. 发送消息
           
           利用websocket发送消息就可以了，但是格式有要求
           
            ```$xslt
             // Example message
             message Event {
             	// unique id
             	string id = 1;
             	// unix timestamp
             	int64 timestamp = 2;
             	// message
             	string message = 3;
             }
            ```
          
            id 需要和组册tcp的时候主题一致，这里是 id="tcp"，message 是具体的消息体，可以是任意格式
            例如：
            ```$xslt
            {"id":"tcp","message":{"name":"zhangsan"}}
            ```
