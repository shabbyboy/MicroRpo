package DB

import (
	"MicroRpo/conf/confserver"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

var (
	DbConfig confserver.Conf
	once sync.Once
	redispool *redis.Pool
	)

func init(){

	DbConfig = confserver.Conf{
		Path:"dbconn/dbconf/db.json",
	}
}

type DbBase interface {
	//获取key
	GetMainKey(gameId string,userId string) string
	//获取一个连接
	NewConn() (redis.Conn,error)
}

type DbConf struct {
	Host string `json:"host"`
	Type string `json:"type"`
	Index []int `json:index`
}

type DbAuth struct {
	Password string `json:"password"`
	MaxIdle int `json:"maxidle"`
	MaxActive int `json:"maxactive"`
	IdleTimeout int `jsong:"idletimeout"`
}

type RedisDB struct {
	FmtKey string
	DbName string
}

func GetDBIndex(dbname string,userId int) int{
	var conf DbConf

	DbConfig.LoacConf()

	DbConfig.ConfExtract(&conf,"database",dbname)

	lenth := len(conf.Index)

	if lenth != 0{
		dbindex := userId%lenth
		return conf.Index[dbindex]
	}
	//返回默认数据库
	return 0
}


// 自己定义的形式和gameId以及userid
func (rb *RedisDB) GetMainKey(gameId int, userId int) string{
	return fmt.Sprintf(rb.FmtKey,gameId,userId)
}

func (rb *RedisDB) NewConn() (redis.Conn,error){

	var (
		dbconf DbConf
		authPass DbAuth
	)

	once.Do(func() {

		DbConfig.LoacConf()
		DbConfig.ConfExtract(&dbconf,"database",rb.DbName)
		DbConfig.ConfExtract(&authPass,"database","auth")

		redispool = &redis.Pool{
			MaxIdle: authPass.MaxIdle,
			MaxActive: authPass.MaxActive,
			IdleTimeout: time.Second * time.Duration(authPass.IdleTimeout),
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(dbconf.Type,dbconf.Host)
				if err != nil{
					return nil,err
				}
				c.Do("auth", authPass.Password)
				return c,err
			},
		}

	})

	conn := redispool.Get()

	//conn.Do("select",dbindex)
	return conn,nil
}

func (rb *RedisDB) SelectDB(conn redis.Conn,indexdb int){
	conn.Do("select",indexdb)
}

func (rb *RedisDB) get(cmd string,gameId int,userId int)(int,error){
	key := rb.GetMainKey(gameId,userId)

	conn,err := rb.NewConn()
	rb.SelectDB(conn,GetDBIndex(rb.DbName,userId))

	defer conn.Close()
	if err != nil{
		return 0,err
	}

	flag, err := redis.Int(conn.Do(cmd,key))
	return flag,err
}

func (rb *RedisDB) randomKey(cmd string,dbindex int)(string,error){
	conn,err := rb.NewConn()
	rb.SelectDB(conn,dbindex)
	defer conn.Close()
	if err != nil{
		return "",err
	}

	flag, err := redis.String(conn.Do(cmd))
	return flag,err
}

func (rb *RedisDB) EXISTS(gameId int,userId int)(int,error){
	return rb.get("EXISTS",gameId,userId)
}

func (rb *RedisDB) DEL(gameId int,userId int)(int,error){
	return rb.get("DEL",gameId,userId)
}

func (rb *RedisDB) RANDOMKEY(dbindex int)(string,error){
	return rb.randomKey("RANDOMKEY",dbindex)
}

func (rb *RedisDB) dbSize(cmd string,dbindex int) (int,error){
	conn,err := rb.NewConn()
	rb.SelectDB(conn,dbindex)
	defer conn.Close()
	if err != nil{
		return -1,err
	}

	flag, err := redis.Int(conn.Do(cmd))
	return flag,err
}

func (rb *RedisDB) DBSIZE(dbindex int)(int,error){
	return rb.dbSize("DBSIZE",dbindex)
}

func (rb *RedisDB) expire(cmd string,gameId int,userId int,time int)(int,error){
	key := rb.GetMainKey(gameId,userId)
	conn,err := rb.NewConn()
	rb.SelectDB(conn,GetDBIndex(rb.DbName,userId))
	defer conn.Close()
	if err != nil{
		return 0,err
	}

	flag, err := redis.Int(conn.Do(cmd,key,time))
	return flag,err
}

func (rb *RedisDB) EXPIRE(gameId int,userId int,time int)(int,error){
	return rb.expire("EXPIRE",gameId,userId,time)
}

func (rb *RedisDB) TTL(gameId int,userId int)(int,error){
	return rb.get("TTL",gameId,userId)
}

func (rb *RedisDB) EXPIREAT(gameId int,userId int,time int)(int,error){
	return rb.expire("EXPIREAT",gameId,userId,time)
}

func (rb *RedisDB) PERSIST(gameId int,userId int)(int,error){
	return rb.get("PERSIST",gameId,userId)
}

func (rb *RedisDB) PEXPIREAT(gameId int,userId int,time int)(int,error){
	return rb.expire("PEXPIREAT",gameId,userId,time)
}

func (rb *RedisDB) PEXPIRE(gameId int,userId int,time int)(int,error){
	return rb.expire("PEXPIRE",gameId,userId,time)
}

func (rb *RedisDB) PTTL(gameId int,userId int)(int,error){
	return rb.get("PTTL",gameId,userId)
}