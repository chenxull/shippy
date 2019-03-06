# 开发记录


## 实现效果

使用微服务的思想完成一个小型的港口货物管理系统


## 问题

### 问题1
 在使用micro-go框架的时候，发现其中许多的包依赖都需要自己手动在github上添加，十分的耗时。需要引入方便的包管理工具
 
### 问题2
 `consignment-cli`和`consignment-service`都成功编译使用容器运行了，在测试的过程中cli始终收不到来自service的数据，在consignment-service容器中已经成功的
 收到了来自cli的数据，目前的问题是cli无法接受到来自service的数据。
 
 **测试**
 在没有使用容器封装二个微服务时，可以进行正常通信
 
 **原因**
 中文代码中存在问题，其定义的接口类型`IRepository`没有使用上，将其改为`Repository`后,并将原来的`Repository`结构体改为了`ConsignmentRepository`
。重新进行编译即可成功运行