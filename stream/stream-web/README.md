### micro web 服务

* http.FileServer(http.Dir("./stream-web/html")) 这里Dir应该是绝对路径 
如果是相对路径的话，需要考虑到go run 和 编译后运行 go 相对路径不同的问题

* 路由的规则，需要和服务名称对应上，路由需要加上/stream 前缀，原因是web handler
处理器代理的原因

* 对原来的服务做了修改，和后端的stream-srv通信，进行了修改


[原代码地址](https://github.com/microhq/stream-web)