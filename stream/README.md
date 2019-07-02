### tcp 服务 
对micro 的一个stream服务用例进行改造，以实现一个tcp 代理

* stream-srv 
    
stream服务，通过绑定路由规则的方式，编写tcp服务

* stream-web

实现websocket提供tcp连接，和客户端进行tcp连接，在通过srv 
clien连接和stream-srv通信，从而完成客户端和srv的通信，作为两者
之间代理