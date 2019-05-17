## handler 为web 类型的网关启动方式

### 先启动srv 服务 

* cd MicroRpo
* go run srv/websrv/main.go

### 启动web服务
* go run web/userweb/main.go

### 启动micro 网关
* micro --api_namespace=microrpo.web api --handler=web
* --api_namespace 指定api网关的命名空间
* --handler 指定处理器类型，处理器类型和命名空间类型需要保持一致

### 测试了下负载均衡效果，发现了下面这个问题
* 再新启动了一个web 服务后，用curl 连续请求，有部分请求会发送到已经关闭的web服务上，需要重启micro api 来解决