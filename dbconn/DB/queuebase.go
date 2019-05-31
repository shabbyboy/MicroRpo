package DB

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type RedisQueue struct {
	RedisDB
}

func (rq *RedisQueue) push(cmd string,gameId int,userId int,data ...interface{}) (int,error){
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}

	args := redis.Args{key}

	for _, v := range data{
		td,_ := json.Marshal(v)
		args = append(args,td)
	}

	ret, err := redis.Int(conn.Do(cmd,args...))

	return ret,err
}

func (rq *RedisQueue) pop(cmd string,gameId int,userId int,ret interface{}) error{
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return err
	}

	data, err := redis.Bytes(conn.Do(cmd,key))

	json.Unmarshal(data,ret)

	return err
}

func (rq *RedisQueue) lrem(cmd string,gameId int,userId int, count int,value interface{}) (int,error){
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}

	ret, err := redis.Int(conn.Do(cmd,key,count,value))

	return ret,err
}

func (rq *RedisQueue) byKeys(cmd string,gameId int,userId int) (int,error){
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}

	ret, err := redis.Int(conn.Do(cmd,key))
	return ret,err
}

func (rq *RedisQueue) lindex(cmd string,gameId int,userId int,index int,data interface{}) error{
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return err
	}

	ret, err := redis.Bytes(conn.Do(cmd,key,index))

	json.Unmarshal(ret,data)

	return err
}

func (rq *RedisQueue) linsert(cmd string,gameId int,userId int,opr string,pos interface{},data interface{}) (int,error){
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return -1,err
	}

	ret, err := redis.Int(conn.Do(cmd,key,opr,pos,data))
	return ret,err
}

func (rq *RedisQueue) lset(cmd string,gameId int,userId int,index int,data interface{}) error{
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return err
	}
	data, _ = json.Marshal(data)

	_, err = redis.Bytes(conn.Do(cmd,key,index,data))

	return err
}

func (rq *RedisQueue) lrange(cmd string,gameId int,userId int,start int,end int) ([][]byte,error){
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return nil,err
	}

	ret, err := redis.ByteSlices(conn.Do(cmd,key,start,end))

	return ret,err
}

func (rq *RedisQueue) ltrim(cmd string,gameId int,userId int,start int,end int) (string,error){
	key := rq.GetMainKey(gameId,userId)

			conn,err := rq.NewConn()
	rq.SelectDB(conn,GetDBIndex(rq.DbName,userId))
	defer conn.Close()
	if err != nil{
		return "",err
	}

	ret, err := redis.String(conn.Do(cmd,key,start,end))

	return ret,err
}



func (rq *RedisQueue) LPUSH(gameId int,userId int,data interface{}) (int,error) {
	return rq.push("LPUSH",gameId,userId,data)
}

func (rq *RedisQueue) LPUSHX(gameId int,userId int,data interface{}) (int,error){
	return rq.push("LPUSHX",gameId,userId,data)
}

func (rq *RedisQueue) RPUSH(gameId int,userId int,data interface{}) (int,error){
	return rq.push("RPUSH",gameId,userId,data)
}

func (rq *RedisQueue) RPUSHX(gameId int,userId int,data interface{}) (int,error){
	return rq.push("RPUSHX",gameId,userId,data)
}

func (rq *RedisQueue) LPOP(gameId int,userId int,ret interface{}) error{
	return rq.pop("LPOP",gameId,userId,ret)
}

func (rq *RedisQueue) RPOP(gameId int,userId int,ret interface{}) error{
	return rq.pop("RPOP",gameId,userId,ret)
}

func (rq *RedisQueue) RPOPLPUSH(){

}

func (rq *RedisQueue) LREM(gameId int,userId int,count int,value interface{})(int,error){
	return rq.lrem("LREM",gameId,userId,count,value)
}

func (rq *RedisQueue) LLEN(gameId int,userId int) (int,error){
	return rq.byKeys("LLEN",gameId,userId)
}

func (rq *RedisQueue) LINDEX(gameId int,userId int,index int,data interface{}) error{
	return rq.lindex("LINDEX",gameId,userId,index,data)
}

func (rq *RedisQueue) LINSERTBEFOR(gameId int,userId int,pos interface{},data interface{}) (int,error){
	return rq.linsert("LINSERT",gameId,userId,"BEFORW",pos,data)
}

func (rq *RedisQueue) LINSERTAFTER(gameId int,userId int,pos interface{},data interface{}) (int,error){
	return rq.linsert("LINSERT",gameId,userId,"AFTER",pos,data)
}

func (rq *RedisQueue) LSET(gameId int,userId int,index int,data interface{}) error{
	return rq.lindex("LSET",gameId,userId,index,data)
}

func (rq *RedisQueue) LRANGE(gameId int,userId int,start int, end int)([][]byte,error){
	return rq.lrange("LRANGE",gameId,userId,start,end)
}

func (rq *RedisQueue) LTRIM(gameId int,userId int,start int, end int)(string,error){
	return rq.ltrim("LTRIM",gameId,userId,start,end)
}