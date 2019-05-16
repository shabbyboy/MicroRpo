package main

import (
	"context"
	"MicroRpo/srv/websrv/msgproto"
	"github.com/micro/go-micro"
	"log"
	"time"
)
/*
写个hello world
 */

type SayHello struct {

}

func (say *SayHello)Hello(ctx context.Context,request *msgproto.Request,response *msgproto.Response) error{
	response.Msg = "hello" + request.Name
	return nil
}

func main(){
	service := micro.NewService(
		micro.Name("microrpo.srv.websrv"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	service.Init()

	msgproto.RegisterSayHelloHandler(service.Server(),new(SayHello))

	if err := service.Run(); err != nil{
		log.Fatal(err)
	}

}
