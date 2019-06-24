package main

import (
	"MicroRpo/rposync/lock"
	"fmt"
	"sync"
	"time"
)

var mutex = sync.Mutex{}
var redislock =lock.RedisLock{}
var sqllock = lock.SQLLock{}
func testRpo(key int,val *string){
	//lock.RpoLock.RequireLock("tlp",val)
	//
	//defer lock.RpoLock.ReleaseLock("tlp",val)

	sqllock.AquireLock()
	defer sqllock.ReleaseLock()

	//mutex.Lock()
	//defer mutex.Unlock()

	//kk := "wangermazi"
	//val = &kk
	fmt.Println(key,val)

}

func main(){
	var temp = "zhangsan"
	t1 := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(num int) {
			testRpo(num,&temp)
			wg.Done()
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(t1)
	fmt.Println("time duration:",elapsed)
	select {

	}
}
