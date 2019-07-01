### tcp 服务

传输数据格式：

```$xslt
type Data struct {
	Cmd string `json:"cmd"`
	Typ SendType `json:"typ"`
	Action string `json:"action"`
	MsgPack map[string]interface{} `json:"msgpack"`
}
```

usage example:

处理消息

```$xslt 
func HandlerHello(id string,options ...interface{}) bool {
	fmt.Println("hello server")
	return true
}
```

注册处理方法

```$xslt
h, err := handler.NewStream()
if err != nil {
    log.Fatal(err)
}
h.Wrapper("abc","",plugins.HandlerHello)
```

发送消息到客户端

```$xslt
//所有的发出去的消息都走这里
//有一条消息体里必须包含一个 typ 必须是SendClient
func SendToClient(msg pb.Message){
	client := client2.DefaultClient

	streamC := pb.NewStreamService("go.micro.srv.stream",client)

	stream, err := streamC.Publish(context.Background())

	if err != nil {
		return
	}

	stream.Send(&msg)
}
```

[源码地址](https://github.com/microhq/stream-srv)
