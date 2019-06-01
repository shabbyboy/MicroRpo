### 封装go-config 简化些代码

* usage

  db.json 配置文件
  
  ```$xslt
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

    1. 默认是静态文件模式，加载本地的静态json.文件

    ```$xslt
    DbConfig := confserver.DefaultConf("dbconn/dbconf/db.json")

    type DbConf struct {
	    Host string `json:"host"`
	    Type string `json:"type"`
	    Index []int `json:index`
    }

    DbConfig.ConfExtract(&conf,"database","user")
    
    

    ```
    2. 新的配置中心，可以是consul 或者 etcd 
    
    ```$xslt
    
    conf := Conf{
        SrcType:CONFIG_CONSUL,
        Coder:CODER_JSON,
        HOST:"10.0.0.10:8500"
        Path:"/my/prefix",
    }
    
    conf.LoadConf()
    
    之后的用法是一样的了

    ```
    
    3. 监听配置文件更新
    ```$xslt
    conf.ConfChange(&host,"database","user")
    ```
    
    
    

