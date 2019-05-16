package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	"log"
	"net/http"
	"MicroRpo/srv/websrv/msgproto"
)

func main(){
	service := web.NewService(
		web.Name("microrpo.web.user"),
	)

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

			writer.Write([]byte(`<html><body><h1>` + rsp.Msg + `</h1></body></html>`))
			return
		}

		fmt.Fprint(writer,`<html><body><h1>Enter Name<h1><form method=post><input name=name type=text /></form></body></html>`)
	})


	if err := service.Init(); err != nil{
		log.Fatal(err)
	}

	if err := service.Run(); err != nil{
		log.Fatal(err)
	}

}

