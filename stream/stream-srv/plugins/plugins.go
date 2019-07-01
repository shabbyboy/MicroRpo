package plugins

import (
	pb "MicroRpo/stream/stream-srv/proto/stream"
	"context"
	"fmt"
	client2 "github.com/micro/go-micro/client"
)

type Plugins func(id string,options ...interface{}) bool

func HandlerHello(id string,options ...interface{}) bool {
	fmt.Println("hello server")
	return true
}


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



