## 注意

* 看了下go-web的代码，发现go-web 启动的是一个http 服务，可以
通过micro.address( ) 指定端口，这样就可以绕过api 网关直接访问web服务，如果不指定端口，
由服务随机指定端口，访问是被拒绝的

* 如下所示： 
> web.Address(":8088") 就可以通过 http://localhost:8088?name=zhansan 直接访问
