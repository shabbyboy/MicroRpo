package DB

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
)

type RedisHash struct {
	RedisDB
}

func (hd *RedisHash) set(cmd string,gameId int,userId int,field string,data interface{}) (string,error){
	key := hd.GetMainKey(gameId,userId)
	conn,err := hd.NewConn(GetDBIndex(hd.DbName,userId))
	defer conn.Close()
	if err != nil{
		return "",err
	}
	data,err = json.Marshal(data)
	if err != nil{
		return "",err
	}
	num, err := redis.String(conn.Do(cmd,key,field,data))
	return num,err
}

func (hd *RedisHash) setMux(cmd string,gameId int,userId int,data map[string]interface{}) (string,error){
	key := hd.GetMainKey(gameId,userId)
	conn,err := hd.NewConn(GetDBIndex(hd.DbName,userId))
	defer conn.Close()
	if err != nil{
		return "",err
	}

	args := redis.Args{key}

	for field,value := range data{
		value,err = json.Marshal(value)
		args = append(args,field,value)
	}
	if err != nil{
		return "",err
	}
	num, err := redis.String(conn.Do(cmd,args...))
	return num,err
}

func (hd *RedisHash) getMux(cmd string,gameId int,userId int,field []string,data ...interface{}) error{
	key := hd.GetMainKey(gameId,userId)
	conn, err := hd.NewConn(GetDBIndex(hd.DbName,userId))
	defer conn.Close()
	if err != nil{
		return err
	}

	args := redis.Args{key}
	for _,v := range field{
		args = append(args,v)
	}
	values, err := redis.ByteSlices(conn.Do(cmd,args...))

	if len(values) != len(data){
		return errors.New("field len not equal data len")
	}

	for k,v := range values{
		json.Unmarshal(v,data[k])
	}
	return err
}

func (hd *RedisHash) getInt(cmd string,gameId int,userId int,field ...string)(int,error){
	key := hd.GetMainKey(gameId,userId)

	conn,err := hd.NewConn(GetDBIndex(hd.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}

	flag, err := redis.Int(conn.Do(cmd,key,field))
	return flag,err
}

func (hd *RedisHash) getFloat(cmd string,gameId int,userId int,field ...string)(float64,error){
	key := hd.GetMainKey(gameId,userId)

	conn,err := hd.NewConn(GetDBIndex(hd.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}

	flag, err := redis.Float64(conn.Do(cmd,key,field))
	return flag,err
}

func (hd *RedisHash) getByKeys(cmd string,gameId int,userId int)([]interface{},error){
	key := hd.GetMainKey(gameId,userId)
	conn,err := hd.NewConn(GetDBIndex(hd.DbName,userId))
	defer conn.Close()
	if err != nil{
		return nil,err
	}
	fieldstr,err := redis.Values(conn.Do(cmd,key))
	return fieldstr,err
}


func (hd *RedisHash) HSET(gameId int,userId int,field string, data interface{}) (string,error){
	num,err := hd.set("HSET",gameId,userId,field,data)
	return num,err
}


func (hd *RedisHash) HSETNX(gameId int,userId int,field string,data interface{})(string, error){
	num,err := hd.set("HSETNX",gameId,userId,field,data)
	return num,err
}


//gameid, userid field data地址
func (hd *RedisHash) HGET(gameId int,userId int,field string,data interface{}) (error){
	key := hd.GetMainKey(gameId,userId)
	conn,err := hd.NewConn(GetDBIndex(hd.DbName,userId))
	//defer conn.Close()
	if err != nil{
		return err
	}
	reData, err := redis.Bytes(conn.Do("HGET",key,field))
	if err != nil{
		return err
	}
	err = json.Unmarshal(reData,data)
	if err != nil{
		return err
	}
	return nil
}

//判断是否存在域
func (hd *RedisHash) HEXISTS(gameId int,userId int,field string) (int,error){
	return hd.getInt("HEXISTS",gameId,userId,field)
}

func (hd *RedisHash) HDEL(gameId int,userId int,field ...string)(int,error){
	return hd.getInt("HDEL",gameId,userId,field...)
}

func (hd *RedisHash) HLEN(gameId int,userId int) (int,error){
	return hd.getInt("HLEN",gameId,userId)
}

func (hd *RedisHash) HSTRLEN(gameId int,userId int,field string)(int,error){
	return hd.getInt("HSTRLEN",gameId,userId,field)
}

func (hd *RedisHash) HINCRBY(gameId int,userId int,field string)(int,error){
	return hd.getInt("HINCRBY",gameId,userId,field)
}

func (hd *RedisHash) HINCRBYFLOAT(gameId int,userId int,field string)(float64,error){
	return hd.getFloat("HINCRBYFLOAT",gameId,userId,field)
}

func (hd *RedisHash) HMSET(gameId int,userId int,data map[string]interface{})(string,error){
	return hd.setMux("HMSET",gameId,userId,data)
}

func (hd *RedisHash) HMGET(gameId int,userId int,field []string,data ...interface{}) error{


	return hd.getMux("HMGET",gameId,userId,field,data...)
}

func (hd *RedisHash) HKEYS(gameId int,userId int)([]interface{},error){
	return hd.getByKeys("HKEYS",gameId,userId)
}

func (hd *RedisHash) HVALS(gameId int,userId int)([]interface{},error){
	return hd.getByKeys("HVALS",gameId,userId)
}

func (hd *RedisHash) HGETALL(gameId int,userId int)([]interface{},error){
	return hd.getByKeys("HGETALL",gameId,userId)
}

