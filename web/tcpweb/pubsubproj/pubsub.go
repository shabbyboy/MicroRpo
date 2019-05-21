package pubsubproj

/*
包装下publiser 和 subscribe
 */

import (
	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
)

type Subscribe struct {
	Server server.Server
}

func (s *Subscribe) SubTopic(topic string,h interface{}) error{

	if err := micro.RegisterSubscriber(topic,s.Server,h); err != nil{
		return err
	}
	return nil
}

func (s *Subscribe) Run() {
	s.Server.Start()
}


type Publish struct {
	ctx context.Context
	client client.Client
}

func (p *Publish) NewPublisher(topic string) micro.Publisher{

	return micro.NewPublisher(topic,p.client)
}


func (p *Publish) PubEvent(publiser micro.Publisher,event interface{}) error{

	if err := publiser.Publish(p.ctx,event);err != nil{
		return err
	}
	return nil
}
