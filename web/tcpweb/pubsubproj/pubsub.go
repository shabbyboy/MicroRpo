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
	Ctx context.Context
	Client client.Client
	//Publishser micro.Publisher
}

func (p *Publish) NewPublisher(topic string) micro.Publisher{

	//p.Publishser = micro.NewPublisher(topic,p.Client)
	return micro.NewPublisher(topic,p.Client)
}


func (p *Publish) PubEvent(publiser micro.Publisher,event interface{}) error{
	//if p.Publishser == nil{
	//	return errors.New("publisher is nil")
	//}

	if err := publiser.Publish(p.Ctx,event);err != nil{
		return err
	}
	return nil
}
