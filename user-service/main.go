package main

import (
	"fmt"
	"log"

	pb "github.com/chenxull/shippy/user-service/proto/user"
	"github.com/micro/go-micro"
)

func main() {

	// 创建数据库连接
	db, err := CreateConnection()
	defer db.Close()

	fmt.Printf("%+v\n", db)
	fmt.Printf("err: %v\n", err)

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	// 自动检查 User 结构是否变化
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}

	tokenService := &TokenService{repo}

	// 创建新的服务
	server := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	server.Init()

	// Register handler
	pb.RegisterUserServiceHandler(server.Server(), &service{repo, tokenService})

	// Run the server
	if err := server.Run(); err != nil {
		fmt.Println(err)
	}
}
