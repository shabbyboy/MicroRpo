package main

import (
	"MicroRpo/dbconn/DB"
	"fmt"
)


func SubAll(channel string,message string){
	fmt.Println("all:",channel,message)
}

func SubOne(channel string,message string){
	fmt.Println("one:",channel,message)
}

func main(){
	subobj := DB.SubRedis{
		RedisDB:DB.RedisDB{
			DbName:"user",
		},
		Quit:make(chan error),
	}

	subobj.Connect()

	subobj.Sub("abc.*",SubAll)
	subobj.Sub("abc.dd",SubOne)

	<- subobj.Quit
}
