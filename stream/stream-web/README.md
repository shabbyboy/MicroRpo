### micro 流式rpc 服务demo(Asim 老哥写的，直接copy过来)

#### 做了点修改，实际操作后发现，asim老哥写的路由哪有些问题

* http.FileServer(http.Dir("./stream-web/html")) 这里Dir应该是绝对路径

* 路由的规则，需要和服务名称对应上，路由需要加上/stream 前缀，原因是web handler处理器代理的原因




[源码地址](https://github.com/microhq/stream-web)