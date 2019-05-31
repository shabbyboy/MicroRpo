package DB

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RedisStr struct {
	RedisDB
}

func (rs *RedisStr) set(cmd string,gameId int,userId int,value string) (string,error){
	key := rs.GetMainKey(gameId,userId)

		conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	defer conn.Close()
	if err != nil{
		fmt.Println("连接错误",err)
		return "0",err
	}
	//data,err := json.Marshal(value)
	flag,err := redis.String(conn.Do(cmd,key,value))

	return flag,err
}

func (rs *RedisStr) setInt(cmd string,gameId int,userId int,num int,value ...string) (int,error){
	key := rs.GetMainKey(gameId,userId)

		conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}
	args := redis.Args{key,num}
	for _,v := range value{
		args = append(args,v)
	}

	//data,err := json.Marshal(value)
	flag,err := redis.Int(conn.Do(cmd,args...))

	return flag,err
}

func (rs *RedisStr) SET(gameId int,userId int,value string) (string,error){
	return rs.set("SET",gameId,userId,value)
}

func (rs *RedisStr) SETNX(gameId int,userId int,value string) (string,error){
	return rs.set("SETNX",gameId,userId,value)
}

func (rs *RedisStr) SETEX(gameId int,userId int,second int,value string)(int,error){
	return rs.setInt("SETEX",gameId,userId,second,value)
}

func (rs *RedisStr) PSETEX(gameId int,userId int,milSecond int,value string)(int,error){
	return rs.setInt("PSETEX",gameId,userId,milSecond,value)
}

func (rs *RedisStr) get(cmd string,gameId int,userId int) (string,error){
	key := rs.GetMainKey(gameId,userId)

		conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	defer conn.Close()
	if err != nil{
		return "",err
	}

	flag,err := redis.String(conn.Do(cmd,key))

	return flag,err
}

func (rs *RedisStr) getset(cmd string,gameId int,userId int,data string) (string,error){
	key := rs.GetMainKey(gameId,userId)

		conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	defer conn.Close()
	if err != nil{
		return "", err
	}
	///data, err = json.Marshal(data)

	res, err := redis.String(conn.Do(cmd,key,data))

	return res,err
}

func (rs *RedisStr) getInt(cmd string,gameId int,userId int,data ...string)(int,error){
	key := rs.GetMainKey(gameId,userId)

	conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}

	args := redis.Args{key}

	for _,v := range data{
		args = append(args, v)
	}

	lenth, err := redis.Int(conn.Do(cmd,args...))
	return lenth,err
}

func (rs *RedisStr) getRange(cmd string,gameId int,userId int,start ...int)(string,error){
	key := rs.GetMainKey(gameId,userId)

	conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	defer conn.Close()
	if err != nil{
		return "",err
	}
	args := redis.Args{key}
	for _,v :=range start{
		args = append(args,v)
	}

	ret, err := redis.String(conn.Do(cmd,args...))
	return ret,err
}

func (rs *RedisStr) incrFloat(cmd string,gameId int,userId int,num ...float64)(float64,error){
	key := rs.GetMainKey(gameId,userId)

		conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}
	args := redis.Args{key}
	for _,v :=range num{
		args = append(args,v)
	}

	ret, err := redis.Float64(conn.Do(cmd,args...))
	return ret,err
}

func (rs *RedisStr) mset(cmd string,gameId int,userId int,data map[string]interface{})(string,error){
	key := rs.GetMainKey(gameId,userId)
	
	conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	
	
	defer conn.Close()
	if err != nil{
		return "",err
	}

	args := redis.Args{key}

	for k,v := range data{
		args = append(args,k,v)
	}

	ret, err := redis.String(conn.Do(cmd,args...))
	return ret,err
}

func (rs *RedisStr) msetNx(cmd string,gameId int,userId int,data map[string]interface{})(int,error){
	key := rs.GetMainKey(gameId,userId)

		conn,err := rs.NewConn()
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}

	args := redis.Args{key}

	for k,v := range data{
		args = append(args,k,v)
	}

	ret, err := redis.Int(conn.Do(cmd,args...))
	return ret,err
}


func (rs *RedisStr) GET(gameId int,userId int) (string,error){
	return rs.get("GET",gameId,userId)
}

func (rs *RedisStr) GETSET(gameId int,userId int,data string) (string,error){
	return rs.getset("GETSET",gameId,userId,data)
}

func (rs *RedisStr) STRLEN(gameId int,userId int)(int, error){
	return rs.getInt("STRLEN",gameId,userId)
}

func (rs *RedisStr) APPEND(gameId int,userId int,appStr string)(int,error){
	return rs.getInt("APPEND",gameId,userId,appStr)
}

func (rs *RedisStr) SETRANGE(gameId int,userId int,offset int,str string)(int,error){
	return rs.setInt("SETRANGE",gameId,userId,offset,str)
}

func (rs *RedisStr) GETRANGE(gameId int,userId int,start int,end int)(string,error){
	return rs.getRange("GETRANGE",gameId,userId,start,end)
}

func (rs *RedisStr) INCR(gameId int,userId int)(int,error){
	return rs.getInt("INCR",gameId,userId)
}

func (rs *RedisStr) INCRBY(gameId int,userId int,incr int)(int,error){
	return rs.setInt("INCRBY",gameId,userId,incr)
}

func (rs *RedisStr) INCRBYFLOAD(gameId int,userId int,incr float64)(float64,error){
	return rs.incrFloat("INCRBY",gameId,userId,incr)
}

func (rs *RedisStr) DECR(gameId int,userId int)(int,error){
	return rs.getInt("DECR",gameId,userId)
}

func (rs *RedisStr) DECRBY(gameId int,userId int,decr int)(int,error){
	return rs.setInt("DECRBY",gameId,userId,decr)
}
