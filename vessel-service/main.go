package main

import (
	"log"
	"os"

	"github.com/micro/go-micro"

	pd "github.com/chenxull/shippy/vessel-service/proto/vessel"
)

const (
	defaultHost = "localhost:27017"
)

// 模拟数据存入数据库
func createDummyData(repo Repository) {
	vessels := []*pd.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	for _, v := range vessels {
		repo.Create(v)
	}
}

func main() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}

	// 创建服务器连接
	session, err := CreateSession(host)
	defer session.Close()
	if err != nil {
		log.Fatalf("Error connecting to datastore: %v", err)
	}

	// 将可供使用的货船提供给VesselRepository创建通信实例的克隆
	repo := &VesselRepository{session.Clone()}

	createDummyData(repo)

	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	server.Init()
	// 有服务和存储实例即可满足提供服务的基本要求
	// 将实现服务端的API注册到服务端
	pd.RegisterVesselServiceHandler(server.Server(), &service{session})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
