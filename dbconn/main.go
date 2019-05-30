package main

import (
	"MicroRpo/dbconn/DB"
	"fmt"
)

func main(){

	hashRedis := DB.RedisHash{
		DB.RedisDB{
			FmtKey:"user:%v:%v",
			DbName:"user",
		},
	}
	res := make(map[string]interface{})
	res["sex"] = 1
	res["height"] = 222
	//目前发现一个问题，存byte 数组取出来的结果不正确了，其他累心没问题
	res["weight"] = []string{"aa","bbb"}

	count, err := hashRedis.Hge(2019,20)

	if err != nil{
		//fmt.Println(err)
	}
	fmt.Println(count)

	var temp map[string]interface{}

	err = hashRedis.HGET(2019,20,"zhangsan",&temp)

	if err != nil{
		fmt.Println(err)
	}

	for k,v := range temp{
		fmt.Println(k,v)
	}

}
