package DB
/*
SSCAN\SINTER\SMOVE\SINTERSTORE\SUNION\SUNIONSTORE\SDIFF\SDIFFSTORE 没有实现，暂时没想好多个key之间的操作怎么弄
 */
import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type RedisSet struct {
	RedisDB
}

func (rs *RedisSet) sadd(cmd string,gameId int,userId int,value ...interface{})(int,error){
	key := rs.GetMainKey(gameId,userId)

	conn,err := rs.NewConn()
	defer conn.Close()
	if err != nil{
		return 0,err
	}
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))

	args := redis.Args{key}
	for _,data :=range value{
		data,_ = json.Marshal(data)
		args = append(args,data)
	}
	ret, err := redis.Int(conn.Do(cmd,args...))

	return ret,err
}

func (rs *RedisSet) spop(cmd string,gameId int,userId int,value interface{}) error{
	key := rs.GetMainKey(gameId,userId)

	conn,err := rs.NewConn()
	defer conn.Close()
	if err != nil{
		return err
	}
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))

	args := redis.Args{key}
	ret, err := redis.Bytes(conn.Do(cmd,args...))

	json.Unmarshal(ret,value)

	return err
}

func (rs *RedisSet) popbyteslice(cmd string,gameId int,userId int,value ...interface{}) ([][]byte,error){
	key := rs.GetMainKey(gameId,userId)

	conn,err := rs.NewConn()
	defer conn.Close()
	if err != nil{
		return nil,err
	}
	rs.SelectDB(conn,GetDBIndex(rs.DbName,userId))

	args := redis.Args{key}
	args = append(args,value...)
	ret, err := redis.ByteSlices(conn.Do(cmd,args...))

	return ret,err
}

func (rs *RedisSet) SADD(gameId int,userId int,value ...interface{})(int,error){
	return rs.sadd("SADD",gameId,userId,value...)
}

func (rs *RedisSet) SISMEMBER(gameId int,userId int,value ...interface{})(int,error){
	return rs.sadd("SISMEMEBER",gameId,userId,value...)
}

func (rs *RedisSet) SPOP(gameId int,userId int,value interface{}) error{
	return rs.spop("SPOP",gameId,userId,value)
}

func (rs *RedisSet) SRANDMEMBER(gameId int,userId int,count int)([][]byte,error){
	return rs.popbyteslice("SRANDMEMBER",gameId,userId,count)
}

func (rs *RedisSet) SREM(gameId int,userId int,value ...interface{})(int,error){
	return rs.sadd("SREM",gameId,userId,value...)
}

func (rs *RedisSet) SCARD(gameId int,userId int)(int,error){
	return rs.sadd("SCARD",gameId,userId)
}

func (rs *RedisSet) SMEMBERS(gameId int,userId int)([][]byte,error){
	return rs.popbyteslice("SMEMBERS",gameId,userId)
}


