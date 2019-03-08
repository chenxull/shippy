package main

import (
	"log"
	"os"

	pd "github.com/chenxull/shippy/consignment-service/proto/consignment"
	vesselProto "github.com/chenxull/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
)

const (
	defaulthost = "localhost:27017"
)

func main() {

	// 获取容器设置的数据库地址环境变量的值
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = defaulthost
	}

	session, err := CreateSession(dbHost)
	// 创建于MongoDB的主会话，需在退出main()时候手动释放连接
	defer session.Close()
	if err != nil {
		log.Fatalf("create session error:%v\n", err)
	}

	// 使用go-micro框架可以简化服务管理
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	// 作为 vessel-service 的客户端
	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", server.Client())

	server.Init()
	//向rRPC服务器注册微服务
	//此时会把我们实现的微服务service与协议中的ShippingServiceServer绑定
	//因为repo的类型为service，service实现了ShippingServiceServer接口，所以可以传入
	pd.RegisterShippingServiceHandler(server.Server(), &handler{session, vesselClient})
	if err := server.Run(); err != nil {
		log.Fatal("failed to serve", err)
	}

}
