package main

import (
	proto "MicroRpo/stream/stream-srv/proto/stream"
	"MicroRpo/stream/stream-web/handler"
	"fmt"
	"github.com/micro/go-log"
	"github.com/micro/go-web"
	"net/http"
	"os"
)

func dir() http.Dir{
	wd,_ := os.Getwd()
	fmt.Println(wd)
	fmt.Println(http.Dir("stream-web/html"))
	return http.Dir("/Users/tugame/newgodemo/microrpo/MicroRpo/stream-web/html")
}

func main() {

	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.stream"),
		web.Version("latest"),
	)



	// setup client
	client := service.Options().Service.Client()

	handler.Client = proto.NewStreamService("go.micro.srv.stream", client)


	// register call handler
	service.HandleFunc("/stream/video", handler.StreamApi)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
