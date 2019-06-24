package lock

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"runtime"
	"strconv"
	"sync"
	"time"
)

/*
基于redis 实现的分布式锁，这里redis 得是单机部署的
 */

func init(){
	conn,_ = redis.Dial("tcp",redishost)
	conn.Do("AUTH","tugame")
}

var conn redis.Conn
const (
	lockedKey string = "lockedkey"
	redishost string = "172.16.8.75:8004"
	defaultstealSecond = 2
	lockcount = 3
	sqlDefaultLockUpdateInterval = time.Second
)
type RedisLock struct{
	sync.RWMutex
	goid uint64
}
//获取goroutine id
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}



//加锁
//可重入
//可抢夺
func (rl *RedisLock) rpoLock(steal bool) (bool,error){

	rl.Lock()
	defer rl.Unlock()
	val := getGID()

	lockval,err := redis.Uint64(conn.Do("GET",lockedKey))

	if err != nil && err != redis.ErrNil{
		return false,err
	}

	if lockval == val {
		return true,nil
	}else {
		if lockval == 0 {
			flag,_ := redis.Uint64(conn.Do("SETNX",lockedKey,val))
			if flag == 1 {
				rl.goid = val
				return true,nil
			}
		}
	}

	return false,nil
}

func (rl *RedisLock) RequireLock(){
	lock,err := rl.rpoLock(false)

	if err != nil {
		panic(err)
	}
	if lock == true {
		return
	} else {
		t := time.AfterFunc(time.Duration(time.Second*defaultstealSecond), func() {
			//一直无法获得锁的情况，判断为死锁
			conn.Do("DEL",lockedKey)
			panic(errors.New("dead lock"))
		})
		for {
			//time.Sleep(time.Duration(1.5 * float64(sqlDefaultLockUpdateInterval)))
			lock,err := rl.rpoLock(false)
			if err != nil{
				fmt.Println(err,1)
				panic(err)
			}
			if lock == true {
				//如果获取到了锁，则停止计时器
				t.Stop()
				return
			}
		}
	}
	}

func (rl *RedisLock) ReleaseLock(){
	rl.Lock()
	defer rl.Unlock()
	lockval,_ := redis.Uint64(conn.Do("GET",lockedKey))

	if rl.goid == lockval {
		conn.Do("DEL",lockedKey)
	}
}
