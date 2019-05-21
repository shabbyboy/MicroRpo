package main

/*
go-config demo
 */

import (
	"fmt"
	"github.com/micro/go-config"
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
	if err := config.LoadFile("conf/config/conf.json");err != nil{
		fmt.Println(err)
	}
	//fmt.Println(os.Getwd())
	//config.load
	//enc := toml.NewEncoder()
	//config.Load(file.NewSource(
	//		file.WithPath("../config/conf.json"),
	//		source.WithEncoder(enc),
	//	))

	var host Host

	//config.Get("hosts","database").Scan(&host)



	w, err := config.Watch("hosts","database")
	if err != nil{
		fmt.Println(err)
	}


	v,err := w.Next()

	if err != nil{
		fmt.Println(err)
	}


	v.Scan(&host)

	fmt.Println(host.Address, host.Port,host.Host)
	//conf := config.Map()
	//for k,v := range conf{
	//
	//	j,_ := v.(map[string]interface{})
	//	for kk,vv := range j{
	//		fmt.Println(kk,vv)
	//	}
	//	fmt.Println(k,v)
	//}
}
