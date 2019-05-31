package DB

import (
	"github.com/garyburd/redigo/redis"
)

type RedisZset struct {
	RedisDB
}

func (rz *RedisZset) zadd(cmd string,gameId int,userId int,value ...interface{}) (int,error){

	key := rz.GetMainKey(gameId,userId)

			conn,err := rz.NewConn()
	rz.SelectDB(conn,GetDBIndex(rz.DbName,userId))
	defer conn.Close()

	if err != nil{
		return 0,err
	}

	args := redis.Args{key}

	args = append(args,value...)

	ret ,err := redis.Int(conn.Do(cmd,args...))

	return ret,err
}

func (rz *RedisZset) zscore(cmd string,gameId int,userId int,value ...interface{}) (string,error){

	key := rz.GetMainKey(gameId,userId)

			conn,err := rz.NewConn()
	rz.SelectDB(conn,GetDBIndex(rz.DbName,userId))
	defer conn.Close()

	if err != nil{
		return "",err
	}

	args := redis.Args{key}

	args = append(args,value...)

	ret ,err := redis.String(conn.Do(cmd,args...))

	return ret,err
}

func (rz *RedisZset) zrange(cmd string,gameId int,userId int,withscore bool,value ...interface{})([][]byte,error){
	key := rz.GetMainKey(gameId,userId)

			conn,err := rz.NewConn()
	rz.SelectDB(conn,GetDBIndex(rz.DbName,userId))
	defer conn.Close()

	if err != nil{
		return nil,err
	}

	args := redis.Args{key}

	args = append(args,value...)

	if withscore{
		args = append(args,"WITHSCORES")
	}
	ret, err := redis.ByteSlices(conn.Do(cmd,args...))

	return ret,err
}

func (rz *RedisZset) zrangebyscore(cmd string,gameId int,userId int,value ...interface{})([][]byte,error){
	key := rz.GetMainKey(gameId,userId)

			conn,err := rz.NewConn()
	rz.SelectDB(conn,GetDBIndex(rz.DbName,userId))
	defer conn.Close()

	if err != nil{
		return nil,err
	}

	args := redis.Args{key}

	args = append(args,value...)

	ret, err := redis.ByteSlices(conn.Do(cmd,args...))

	return ret,err
}

//没想到好的添加方法
func (rz *RedisZset) ZADD(gameId int,userId int,value ...interface{}) (int,error){
	return rz.zadd("ZADD",gameId,userId,value...)
}

func (rz *RedisZset) ZSCORE(gameId int,userId int,value interface{})(string,error){
	return rz.zscore("ZSCORE",gameId,userId,value)
}

func (rz *RedisZset) ZINCRBY(gameId int,userId int,value ...interface{})(string,error){
	return rz.zscore("ZINCRBY",gameId,userId,value...)
}

func (rz *RedisZset) ZCARD(gameId int,userId int)(int,error){
	return rz.zadd("ZCARD",gameId,userId)
}

func (rz *RedisZset) ZCOUNT(gameId int,userId int,value ...interface{})(int,error){
	return rz.zadd("ZCOUNT",gameId,userId,value...)
}

func (rz *RedisZset) ZRANGEWITHSCORE(gameId int,userId int,withScore bool,value ...interface{})([][]byte,error){
	return rz.zrange("ZRANGE",gameId,userId,true,value...)
}

func (rz *RedisZset) ZRANGE(gameId int,userId int,withScore bool,value ...interface{})([][]byte,error){
	return rz.zrange("ZRANGE",gameId,userId,false,value...)
}

func (rz *RedisZset) ZREVRANGEWITHSCORE(gameId int,userId int,withScore bool,value ...interface{})([][]byte,error){
	return rz.zrange("ZREVRANGE",gameId,userId,true,value...)
}

func (rz *RedisZset) ZREVRANGE(gameId int,userId int,withScore bool,value ...interface{})([][]byte,error){
	return rz.zrange("ZREVRANGE",gameId,userId,false,value...)
}

func (rz *RedisZset) ZRANGEBYSCORE(gameId int,userId int,value ...interface{})([][]byte,error){
	return rz.zrangebyscore("ZRANGEBYSCORRE",gameId,userId,value...)
}

func (rz *RedisZset) ZREVRANGEBYSCORE(gameId int,userId int,value ...interface{})([][]byte,error){
	return rz.zrangebyscore("ZREVRANGEBYSCORRE",gameId,userId,value...)
}

func (rz *RedisZset) ZRANK(gameId int,userId int,value interface{})(int,error){
	return rz.zadd("ZRANK",gameId,userId,value)
}

func (rz *RedisZset) ZREVRANK(gameId int,userId int,value interface{})(int,error){
	return rz.zadd("ZREVRANK",gameId,userId,value)
}

func (rz *RedisZset) ZREM(gameId int,userId int,value ...interface{})(int,error){
	return rz.zadd("ZREM",gameId,userId,value)
}
//根据索引移除 和这个名字不相符啊
func (rz *RedisZset) ZREMRANGEBYRANK(gameId int,userId int,value ...interface{}) (int,error){
	return rz.zadd("ZREMRANGEBYRANK",gameId,userId,value...)
}

func (rz *RedisZset) ZREMRANGEBYSCORE(gameId int,userId int,value ...interface{})(int, error){
	return rz.zadd("ZREMRANGEBYSCORE",gameId,userId,value...)
}

func (rz *RedisZset) ZRANGEBYLEX(gameId int,userId int,value ...interface{}) ([][]byte,error){
	return rz.zrangebyscore("ZRANGEBYLEX",gameId,userId,value...)
}

func (rz *RedisZset) ZLEXCOUNT(gameId int,userId int,value ...interface{}) (int,error){
	return rz.zadd("ZLEXCOUNT",gameId,userId,value...)
}

func (rz *RedisZset) ZREMRANGEBYLEX(gameId int,userId int,value ...interface{}) (int,error){
	return rz.zadd("ZREMRANGEBYLEX",gameId,userId,value...)
}

func (rz *RedisZset) ZUNIONSTORE(gameId int,userId int,value ...interface{})(int,error){
	return rz.zadd("ZUNIONSTORE",gameId,userId,value...)
}

func (rz *RedisZset) ZINTERSTORE(gameId int,userId int,value ...interface{})(int,error){
	return rz.zadd("ZINTERSTORE",gameId,userId,value...)
}

