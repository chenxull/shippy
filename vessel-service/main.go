package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/micro/go-micro"

	pd "github.com/chenxull/shippy/vessel-service/proto/vessel"
)

type Repository interface {
	FindAvailable(*pd.Specification) (*pd.Vessel, error)
}

type VesselRepository struct {
	vessels []*pd.Vessel
}

// 实现接口
func (repo *VesselRepository) FindAvailable(spec *pd.Specification) (*pd.Vessel, error) {
	for _, v := range repo.vessels {
		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight {
			return v, nil
		}
	}
	return nil, errors.New("No vessel can`t be use")
}

// 定义货船服务
type service struct {
	repo Repository
}

// 实现服务端
func (s *service) FindAvailable(ctx context.Context, spec *pd.Specification, resp *pd.Response) error {
	// 调用内部方法查找是否有可用的货船
	v, err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}
	fmt.Println("I HAVA GET THE MESSAGES!")
	resp.Vessel = v
	return nil
}

func main() {
	// 停留在港口的货船
	vessels := []*pd.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	// 将可供使用的货船提供给VesselRepository创建实例repo
	repo := &VesselRepository{vessels}
	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	server.Init()
	// 有服务和存储实例即可满足提供服务的基本要求
	// 将实现服务端的API注册到服务端
	pd.RegisterVesselServiceHandler(server.Server(), &service{repo})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
