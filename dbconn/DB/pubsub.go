package DB

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"regexp"
	"strings"
	"unsafe"
)

type SubCallBack func(channel,message string)


type SubRedis struct {
	RedisDB
	client redis.PubSubConn
	callMap map[string]SubCallBack
	//对外传递订阅关闭的信号
	Quit chan error
}

func (ps *SubRedis) Close(){
	err := ps.client.Close()
	if err != nil{
		fmt.Println(err)
	}
	ps.Quit <- errors.New("quit")
}

func (ps *SubRedis) Sub(channel string,scb SubCallBack) error{

	var err error

	if strings.Contains(channel,"*"){
		err = ps.client.PSubscribe(channel)
	}else {
		err = ps.client.Subscribe(channel)
	}

	if err != nil{
		return err
	}
	ps.callMap[channel] = scb
	return nil
}

func (ps *SubRedis) UnSub(channel string) error{
	_, ok := ps.callMap[channel]

	if ok{
		var err error
		if strings.Contains(channel,"*"){
			err = ps.client.PUnsubscribe(channel)
		}else {
			err = ps.client.Unsubscribe(channel)
		}
		if err != nil{
			return err
		}
		delete(ps.callMap, channel)
	}else{
		return errors.New("have not this sub")
	}
	return nil
}

func (ps *SubRedis) Connect() error{
	conn,err := ps.NewConn()

	if err != nil{
		return err
	}

	ps.Quit = make(chan error)
	ps.client = redis.PubSubConn{conn}
	ps.callMap = make(map[string]SubCallBack)

	go func() {

		for{
			switch res := ps.client.Receive().(type) {
			/*
			发现一个很好玩的是PMessage 包含Message，Message 对于通配符的主题无法接收到，所以用PMessage吧，所有主题都可以接收到
			 */
			case redis.PMessage:
				channel := (*string)(unsafe.Pointer(&res.Channel))
				message := (*string)(unsafe.Pointer(&res.Channel))
				fmt.Println("走到message 了:",*channel,*message)
				for ch,call := range ps.callMap{
					match, _ := regexp.MatchString(ch,*channel)
					if match{
						fmt.Println("debug:",ch)
						call(*channel,*message)
					}
				}
			//case redis.Message:
			//	channel := (*string)(unsafe.Pointer(&res.Channel))
			//	message := (*string)(unsafe.Pointer(&res.Channel))
			//	fmt.Println("能走到pmessage 吗",*channel,*message)
			case redis.Subscription:
			case error:
				continue
			}
		}
	}()

	return nil
}


type PubRedis struct {
	RedisDB
}

func (pr *PubRedis) Pub(channel string,data interface{}) error{
	conn,err := pr.NewConn()
	defer conn.Close()
	_,err = conn.Do("PUBLISH",channel,data)

	if err != nil{
		return err
	}
	return nil
}




