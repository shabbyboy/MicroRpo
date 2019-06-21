package lock

import (
	"errors"
	"sync"
)


func init(){
	RpoLock = rpoLock{
		Mutex:sync.Mutex{},
	}
}

var RpoLock rpoLock

/*
基于进程的锁

效果不太理想，我也没想明白这个锁存在的理由，
先放这吧，再看看

 */
type rpoLock struct {
	sync.Mutex
	mugroup map[interface{}]interface{}
}

func (rl *rpoLock) requireLock(lockKey,lockVal interface{}) bool{
	rl.Lock()

	defer rl.Unlock()

	_,ok := rl.mugroup[lockKey]

	if ok {
			return false
	}else {
		if rl.mugroup == nil {
			rl.mugroup = make(map[interface{}]interface{})
		}
		rl.mugroup[lockKey] = lockVal
		return true
	}
}

func (rl *rpoLock) RequireLock(lockKey,lockVal interface{}) {

	for {
		flag := rl.requireLock(lockKey,lockVal)

		if flag {
			return
		}
	}

}

func (rl *rpoLock) releaseLock(lockKey,lockval interface{}) bool {
	rl.Lock()
	defer rl.Unlock()
	val, ok := rl.mugroup[lockKey]

	if ok {
		if val == lockval {
			delete(rl.mugroup,lockKey)
			return true
		}else {
			panic(errors.New("2unlock of unlocked mutex"))
		}
	}else {
		panic(errors.New("1unlock of unlocked mutex"))
	}
}

func (rl *rpoLock) ReleaseLock(lockKey,lockval interface{}) {
	rl.releaseLock(lockKey,lockval)
}


