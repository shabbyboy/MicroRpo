package main

/*
go-config demo
 */

import (
	"MicroRpo/conf/confserver"
	"fmt"
	"time"
)

type Host struct {
	Address string `json:"address"`
	Port int `json:"port"`
	Host map[string]interface{} `json:"host"`
}

type Config struct {
	Hosts Host `json:"hosts"`
}

func main(){
	//conf := config.NewConfig()
	//这个路径和我常识不太一样


	conf := confserver.Conf{
		"conf/config/conf.json",
	}

	if err := conf.LoacConf();err != nil{
		fmt.Println(err)
	}

	var host Host

	conf.ConfExtract(&host,"hosts","database")
	fmt.Println(host.Port,host.Address)

	conf.ConfChange(&host,"hosts","database")

	timeticker := time.NewTicker(time.Second*2)

	for{

		select {
			case <- timeticker.C:
				fmt.Println(host.Address,host.Port)
		}

	}
}
