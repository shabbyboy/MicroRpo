package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	proto "MicroRpo/stream/stream-srv/proto/stream"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	Client proto.StreamService
)

func videoClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// subscribe to the stream
	stream, err := Client.Subscribe(context.TODO(), &proto.SubscribeRequest{
		Id: r.Form.Get("id"),
	})

	if err != nil {
		fmt.Println("Subscribe error for", r.Form.Get("id"), err.Error())
		return
	}

	for {
		msg, err := stream.Recv()
		// stream ended
		if err == io.EOF {
			return
		}
		// some other error
		if err != nil {
			fmt.Println("Stream receive error for", r.Form.Get("id"), err.Error())
			return
		}
		// write the data to the websocket
		if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
			fmt.Println("Write message error for", r.Form.Get("id"), err.Error())
			return
		}
	}
}

func StreamVideo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")
	typ := r.Form.Get("type")
	fmt.Println("daozhel e")

	if len(id) == 0 {
		http.Error(w, "id not set", 500)
		return
	}

	// client
	if typ == "client" {
		fmt.Println("Subscribing to stream", id)
		videoClient(w, r)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println("Publishing stream", id)

	// create a video stream
	_, err = Client.Create(context.TODO(), &proto.CreateRequest{
		Id: id,
	})
	if err != nil {
		fmt.Println("Error creating stream", id, err.Error())
		return
	}

	// send stream
	stream, err := Client.Publish(context.TODO())
	if err != nil {
		fmt.Println("Error publishing stream", id, err.Error())
		return
	}

	// send loop
	for {
		// read from websocket
		_, d, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println("Error reading message", id, err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		// send to server
		if err := stream.Send(&proto.Message{
			Id:   id,
			Data: d,
		}); err != nil {
			fmt.Println("Error sending message", id, err.Error())
			return
		}
	}
}
