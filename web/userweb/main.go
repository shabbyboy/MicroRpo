package main

/*
利用goweb 模块弄了一个web服务
 */

import (
	"MicroRpo/srv/websrv/msgproto"
	"MicroRpo/web/tcpweb/pubsubproj"
	"MicroRpo/web/tcpweb/tcpproto"
	"context"
	"fmt"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	"log"
	"net/http"
)

type pub struct {


}

func (p *pub) Process(ctx context.Context,event *pubsub.Event) error{

	//// 通过发布订阅，往长连接发送消息
	pubproj := pubsubproj.Publish{
		Ctx:context.Background(),
		Client:client.DefaultClient,
	}

	publisher := pubproj.NewPublisher("123")

	if err := pubproj.PubEvent(publisher,event); err != nil{
		fmt.Println("发布：",err)
	}
	return nil
}


func main(){
	service := web.NewService(
		web.Name("microrpo.web.user"),
	)
	//组册订阅，获取长连接的消息

	subproj := pubsubproj.Subscribe{
		//Server:server.DefaultServer,
	}
	err := subproj.SubTopic("tcp",new(pub))

	if err != nil{
		log.Fatal(err)
	}
	subproj.Run()
	//
	//micro.RegisterSubscriber("tcp",server.DefaultServer,new(pub))
	//server.DefaultServer.Start()

	service.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
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


	if err := service.Init(); err != nil{
		log.Fatal(err)
	}

	//service.Options().Service.Run()

	if err := service.Run(); err != nil{
		log.Fatal(err)
	}

}

