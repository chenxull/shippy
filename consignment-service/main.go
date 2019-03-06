package main

import (
	"context"
	"fmt"
	"log"

	pd "github.com/chenxull/shippy/consignment-service/proto/consignment"
	"github.com/micro/go-micro"
)

const (
	PORT = ":50051"
)

type Repository interface {
	Create(*pd.Consignment) (*pd.Consignment, error) //存放新货物
	GetALL() []*pd.Consignment                       //获取仓库中的所有货物
}

//存放多批货物的仓库实现了IRepository接口
type ConsignmentRepository struct {
	consignments []*pd.Consignment
}

func (repo *ConsignmentRepository) Create(consignment *pd.Consignment) (*pd.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *ConsignmentRepository) GetALL() []*pd.Consignment {
	return repo.consignments
}

//定义微服务
type service struct {
	repo Repository
}

//service 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端

//托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pd.Consignment, resp *pd.Response) error {
	//接受承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	fmt.Printf("DEBUG")
	//resp = &pd.Response{Created: true, Consignment: consignment}
	resp.Created = true
	resp.Consignment = consignment
	fmt.Println(resp)
	return nil
}

// 获取目前所托运的货物
func (s *service) GetConsignment(ctx context.Context, req *pd.GetRequest, resp *pd.Response) error {
	allconsignments := s.repo.GetALL()
	//resp = &pd.Response{Consignments: allconsignments}
	resp.Consignments = allconsignments
	return nil
}

func main() {

	// 使用go-micro框架可以简化服务管理
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	server.Init()
	repo := &ConsignmentRepository{}
	//向rRPC服务器注册微服务
	//此时会把我们实现的微服务service与协议中的ShippingServiceServer绑定
	//因为repo的类型为service，service实现了ShippingServiceServer接口，所以可以传入
	pd.RegisterShippingServiceHandler(server.Server(), &service{repo})
	if err := server.Run(); err != nil {
		log.Fatal("failed to serve", err)
	}

}
