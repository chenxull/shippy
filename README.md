# 货运微服务系统
## 最终效果
使用微服务的思想完成一个小型的港口货物管理系统

## 开发记录
### 2019.3.7
使用go-micro 替代了 gRPC增;加货运服务，目前提供三个微服务：
1. consignment-cli
2. consignment-service
3. vessel-service

上述三个微服务的通信逻辑如下图所示：
![](https://images.yinzige.com/2018-05-22-094448.png)


### 2019.3.8
1. 对`consignment-service`和`vessel-service`模块的代码进行重构，将业务逻辑和服务交互逻辑从主文件中提取出来，使得代码结构更加符合实际
2. 增加`user-service`和`user-cli`
3. 使用 GORM 库与 Postgres 数据库进行交互，并将命令行的数据存储进去
4. 使用docker-compose对各个组件进行统一管理
## 运行

每个微服务都使用`docker`进行封装，`Makefile`进行编译。可以通过如下步骤启动整个服务：
1. `cd /vessel-service`--->`make build && make run`
2. `cd /consignment-service`--->`make build && make run`
3. `cd /consignment-cli`--->`make build && make run`




## Problem 

### NO.1
 在使用micro-go框架的时候，发现其中许多的包依赖都需要自己手动在github上添加，十分的耗时。需要引入方便的包管理工具
 
### No.2
 `consignment-cli`和`consignment-service`都成功编译使用容器运行了，在测试的过程中cli始终收不到来自service的数据，在consignment-service容器中已经成功的
 收到了来自cli的数据，目前的问题是cli无法接受到来自service的数据。
 
 **测试**
 在没有使用容器封装二个微服务时，可以进行正常通信
 
 **原因**
 中文代码中存在问题，其定义的接口类型`IRepository`没有使用上，将其改为`Repository`后,并将原来的`Repository`结构体改为了`ConsignmentRepository`
。重新进行编译即可成功运行

### NO.3 