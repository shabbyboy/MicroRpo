package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"MicroRpo/stream/stream-srv/handler"

	pb "MicroRpo/stream-srv/proto/stream"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.stream"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	h, err := handler.NewStream()
	if err != nil {
		log.Fatal(err)
	}

	// Register Handler
	pb.RegisterStreamHandler(service.Server(), h)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
