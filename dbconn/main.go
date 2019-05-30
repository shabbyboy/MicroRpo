package main

import (
	"MicroRpo/dbconn/DB"
	"fmt"
)

func main(){

	queue := DB.RedisQueue{
		DB.RedisDB{
			"queue:%v:%v",
			"user",
		},
	}

	dict := make(map[string]int)

	dict["zhangsan"] = 12
	dict["lisi"] = 15

	queue.PUSH(2019,30,dict)

	temp, err := queue.POP(2019,30)

	if err != nil{
		fmt.Println(err)
	}

	tm,ok := temp.(map[string]int)
	fmt.Println(temp)
	if ok {
		fmt.Println(tm)
	}else {
		fmt.Println("temp not dit" ,tm)
	}
}
