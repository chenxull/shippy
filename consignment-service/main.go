package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pd "github.com/chenxull/shippy/consignment-service/proto/consignment"
)

const (
	PORT = ":50051"
)

type IRepository interface {
	Create(consignment *pd.Consignment) (*pd.Consignment, error) //存放新货物
	GetALL() []*pd.Consignment                                   //获取仓库中的所有货物
}

//存放多批货物的仓库实现了IRepository接口
type Repository struct {
	consignments []*pd.Consignment
}

func (repo *Repository) Create(consignment *pd.Consignment) (*pd.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetALL() []*pd.Consignment {
	return repo.consignments
}

//定义微服务
type service struct {
	repo Repository
}

//service 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端

//托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pd.Consignment) (*pd.Response, error) {
	//接受承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	resp := &pd.Response{Created: true, Consignment: consignment}
	return resp, nil
}

// 获取目前所托运的货物
func (s *service) GetConsignment(context.Context, *pd.GetRequest) (*pd.Response, error) {
	allconsignments := s.repo.GetALL()
	resp := &pd.Response{Consignments: allconsignments}
	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("failed to listen:%v", err)
	}
	log.Printf("listen on : %s \n", PORT)

	server := grpc.NewServer()
	repo := Repository{}

	//向rRPC服务器注册微服务
	//此时会把我们实现的微服务service与协议中的ShippingServiceServer绑定
	//因为repo的类型为service，service实现了ShippingServiceServer接口，所以可以传入
	pd.RegisterShippingServiceServer(server, &service{repo})
	if err := server.Serve(listener); err != nil {
		log.Fatal("failed to serve:%v", err)
	}

}
