package main

/*
此main 函数 主要是开发时 测试用，并无其他用处，可以删除
 */

import (
	"MicroRpo/srv/websrv/msgproto"
	"MicroRpo/web/tcpweb/pubsubproj"
	"MicroRpo/web/tcpweb/tcpproto"
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-web"
	"log"
	"net/http"
)

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


func main(){
	service := web.NewService(
		web.Name("microrpo.web.user"),
	)
	//组册订阅，获取长连接的消息

	if err := service.Init(); err != nil{
		log.Fatal(err)
	}

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

	micro.RegisterSubscriber("tcp",server.DefaultServer,new(pub))
	server.DefaultServer.Start()

	//完整的路由路径的是/user/ 开头，和web 代理有一定的关系
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

}

