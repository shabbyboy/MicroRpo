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