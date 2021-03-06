### redis 数据库封装

* usage
 
 配置文件 db.json
 
```  
{
 "database": {
   "user": {
     "host": "172.16.8.75:8003",
     "type": "tcp",
     "index": [7,8]
   },
   "default": {
     "host": "172.16.8.75:8003",
     "type": "tcp",
     "index": [1,2]
   },
   "auth": {
     "password": "tugame",
     "maxidle": 1,
     "maxactive": 500,
     "idletimeout": 20
   }
 }
}
 ```
 
user和default 是数据库名，实际使用过程中，通过数据库名找到配置信息，index 是具体的数据库索引，
key按照一定的算法会被分配到index 属性中的某一个数据库

 example

```$xslt
hashRedis := DB.RedisHash{
    DB.RedisDB{
        FmtKey:"user:%v:%v",
        DbName:"user",
    },
}
err = hashRedis.HGET(2019,20,"zhangsan",&temp)
```

FmtKey 是key的格式化字符串，dbname 对应配置文件中的数据库


### pubsub usage

* sub topic 订阅的写法

```$xslt

func SubAll(channel string,message string){
	fmt.Println("all:",channel,message)
}
subobj := DB.SubRedis{
    RedisDB:DB.RedisDB{
        DbName:"user",
    },
    Quit:make(chan error),
}

subobj.Connect()

subobj.Sub("abc.*",SubAll)
subobj.Sub("abc.dd",SubOne)

<- subobj.Quit
```

* pub topic 发布的写法

```$xslt
pub := DB.PubRedis{
    RedisDB:DB.RedisDB{
        DbName:"user",
    },
}

err := pub.Pub("abc.ZZZ","wangbadan")
fmt.Println("zzz:",err)

err = pub.Pub("abc.dd","nihao")
fmt.Println("dd:",err)
```