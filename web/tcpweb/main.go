package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/micro/go-web"
	"log"
	"net/http"
	"time"
)

type tcpHandler struct {

}

var(
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)


func (tcp tcpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ws,err := upgrader.Upgrade(w,r,nil)

	if err != nil{
		return
	}

	fmt.Println("daozhele")

	go func() {
		for{
			senderr := ws.WriteMessage(websocket.TextMessage ,[]byte("nihao"))
			if senderr != nil{
				fmt.Println("发送数据出错",senderr)
				ws.Close()
				return
			}
			time.Sleep(1*time.Second)

		}
	}()

	var (
		data []byte
		readerr error
		lenth int
	)
	for {
		if lenth,data,readerr = ws.ReadMessage();readerr != nil{
			ws.Close()
		}
		if readerr = ws.WriteMessage(lenth,data);readerr !=nil{
			ws.Close()
		}
	}

	//fmt.Println(r.Body)

}

func main(){


	wsserver := web.NewService(
		web.Name("microrpo.web.wstcp"),
		//web.Address(":8088"),
	)

	if err := wsserver.Init(); err != nil{
		log.Fatal(err)
	}
	wsserver.Handle("/",new(tcpHandler))


	if err := wsserver.Run(); err != nil{
		log.Fatal(err)
	}

}
