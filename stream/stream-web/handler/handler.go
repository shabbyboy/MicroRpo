package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	proto "MicroRpo/stream/stream-srv/proto/stream"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	Client proto.StreamService
)

//发送到客户端的数据
func StreamToClient(conn *websocket.Conn, id string, ch chan struct{}) {
	// subscribe to the stream
	stream, err := Client.Subscribe(context.TODO(), &proto.SubscribeRequest{
		Id: id,
	})

	if err != nil {
		fmt.Println("Subscribe error for", id, err.Error())
		goto out
	}

	for {
		msg, err := stream.Recv()
		// stream ended
		if err == io.EOF {
			break
		}
		// some other error
		if err != nil {
			fmt.Println("Stream receive error for", id, err.Error())
			break
		}
		// write the data to the websocket
		if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
			fmt.Println("Write message error for", id, err.Error())
			break
		}
	}
	out:
		close(ch)
		ch = nil
}

//发送到服务端的数据
func StreamToServer(conn *websocket.Conn,id string,ch chan struct{}) {

	// create a video stream
	_, err := Client.Create(context.TODO(), &proto.CreateRequest{
		Id: id,
	})
	var stream proto.Stream_PublishService

	if err != nil {
		fmt.Println("Error creating stream", id, err.Error())
		/*
		发现一个很有意思的问题，goto 不能使用在变量声明的前面
		是不是goto 不能跳过变量的声明？
		 */
		goto out
	}
	// send stream
	stream, err = Client.Publish(context.TODO())

	if err != nil {
		fmt.Println("Error publishing stream", id, err.Error())
		goto out
	}

	// send loop
	for {
		// read from websocket
		_, d, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// send to server
		if err := stream.Send(&proto.Message{
			Id:   id,
			Data: d,
		}); err != nil {
			fmt.Println("Error sending message", id, err.Error())
			break
		}
	}
out:
	close(ch)
	ch = nil

}

func StreamApi(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	//换个思路，id 指代游戏id
	//第二天 换个思路也不行啊 😢
	id := r.Form.Get("id")
	fmt.Println("daozhel e")

	if len(id) == 0 {
		http.Error(w, "id not set", 500)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	over := make(chan struct{})

	go StreamToServer(conn,id,over)

	go StreamToClient(conn,id,over)

	<- over
}
