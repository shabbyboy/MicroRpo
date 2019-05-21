### go-config 简单封装

* LoacConf 加载配置
    * 每次使用必须现调用此方法，加载配置
* ConfChange 开启一个协程获取配置文件变更
    * 用于动态配置，传入需要解析成的结构体，和监视的配置key
* ConfExtract 获取配置 
    * 获取指定key 的 value
