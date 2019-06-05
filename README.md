### 目录

* 配置 模块：
    1. [配置支持consul、etcd、静态文件作为配置中心](https://github.com/shabbyboy/MicroRpo/tree/master/conf)

* redis数据库模块：
    1. [对redgo 进行了封装，简化了redis的使用](https://github.com/shabbyboy/MicroRpo/tree/master/dbconn)
    
* 长连接和短连接模块
    1. [利用websocket实现了长连接，短连接沿用了micro的web 处理模式](https://github.com/shabbyboy/MicroRpo/tree/master/web/userweb)
    
* 日志模块
    1. [日志模块用的是logrus框架，要问原因，✨最多](https://github.com/shabbyboy/MicroRpo/tree/master/rpolog)
    
    
* 项目的bin目录用于存放执行文件
    
    1. 格式为：bin/包名/执行文件 
    2. 可以使用提供的build.sh 脚本编译代码
        
        执行脚本
        > sh build.sh
        
        Enter build path: 具体路径，这里是我自己的bin目录
        > Enter build path:/Users/tugame/newgodemo/microrpo/MicroRpo/bin 
        
        Enter code path: 代码路径，main.go 所在目录
        > Enter code path:/Users/tugame/newgodemo/microrpo/MicroRpo/web/userweb
    
    完成之后就可以在bin目录下看到 /userweb/userweb 的执行文件

* runlogs目录是日志输出目录，
    1. 格式：runlogs/包名/日志文件 例如：runlogs/userweb/web23967.log.2019-06-03
