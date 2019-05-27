package main

import (
	"MicroRpo/web/tcpweb/pubsubproj"
	"MicroRpo/web/tcpweb/tcpproto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-web"
	"log"
	"net/http"
	"time"
)

type tcpHandler struct {
	//inchan chan[]byte
}



var(
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)


type wsSub struct {
	inchan chan[]byte
}

func (p *wsSub) Process(ctx context.Context,event *pubsub.Event) error{
	fmt.Println("订阅主题")
	fmt.Println(event.Message)
	if len(event.Message) > 0 {
		fmt.Println(event.Message)
		p.inchan <- []byte(event.Message)
	}
	return nil
}


type wsConnect struct {
	//标识 websocket
	wsId string
	//收到的消息
	inchan chan []byte
	//发送的消息
	outchan chan []byte
	//websocket连接
	wsconn *websocket.Conn
	//终止链路
	stop chan error
}

//发送心跳的逻辑
func (ws wsConnect) SendHeart(data []byte) error{
	//心跳逻辑
	go func() {
		t := time.NewTicker(time.Second*1)
		for{
			select{
			case <- t.C:
				hearterr := ws.wsconn.WriteMessage(websocket.TextMessage ,data)
				if hearterr != nil{
					ws.stop<- hearterr
					return
				}
			}
		}
	}()

	return nil
}


func (ws wsConnect) ReadLoop() error{
	go func() {
		for {
			_, msg, err := ws.wsconn.ReadMessage()
			if err != nil {
				ws.stop <- err
			}
			//ws.inchan <- msg
			ev := &pubsub.Event{
			}


			seqjsonerr := json.Unmarshal(msg,ev)

			if seqjsonerr != nil{
				fmt.Println(seqjsonerr)
			}

			pubproj := pubsubproj.Publish{
				Ctx:context.Background(),
				Client:client.DefaultClient,
			}

			publisher := pubproj.NewPublisher(ev.Id)
			fmt.Println(string(ev.Id))
			if errpub := pubproj.PubEvent(publisher,ev); errpub != nil{
				log.Println(errpub)
			}
		}
	}()
	return nil
}



func (ws wsConnect) WriteLoop() error{

	go func() {

		for {
			data := <- ws.outchan
			if err := ws.wsconn.WriteMessage(websocket.TextMessage,data); err != nil{
				log.Println(err)
				ws.stop <- err
			}
		}
	}()

	return nil
}



func (tcp tcpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		r.ParseForm()

		userid := r.Form.Get("userid")

		ws,err := upgrader.Upgrade(w,r,nil)

		if err != nil{
			return
		}

		wsconn := wsConnect{
			//这个连接id 暂时没想好怎么写，应该更具业务去判断，能够唯一标示客户端
			wsId:userid,
			inchan:make(chan []byte,2),
			outchan:make(chan []byte,2),
			wsconn:ws,
			stop:make(chan error,2),
		}

		//心跳逻辑
		//date := time.Now()
		//wsconn.SendHeart([]byte(date.Format("2006/01/02")+"heartbeat"))

		wsconn.ReadLoop()
		//根据websocket连接id，注册一个发布订阅
		wsconn.WriteLoop()
		tcpHand := &wsSub{inchan:make(chan []byte,2)}
		//server.Init()

		subproj := pubsubproj.Subscribe{
			server.DefaultServer,
		}
		fmt.Println(wsconn.wsId)
		//该用包装后的注册服务
		errregis := subproj.SubTopic(wsconn.wsId,tcpHand)
		subproj.Run()
		///errregis := micro.RegisterSubscriber(wsconn.wsId,server.DefaultServer,tcpHand)
		//server.DefaultServer.Start()

		if errregis != nil {
			fmt.Println(errregis)
		}

		go func() {
			for {
				result := <- tcpHand.inchan
				wsconn.outchan <- result
			}
		}()

		//阻塞直到 websocket 关闭
		<- wsconn.stop
		wsconn.wsconn.Close()
		//这里手动关下rpc服务把避免开启太多
		server.DefaultServer.Stop()
		//log.Fatal("wbsocket 连接关闭")
	}


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

	//server.DefaultServer.Start()
	//wsserver.Options().Service.Run()

	if err := wsserver.Run(); err != nil{
		log.Fatal(err)
	}

}
