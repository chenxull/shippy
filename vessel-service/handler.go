package main

/* 

type VesselServiceHandler interface {
	// 检查是否有能运送货物的轮船
	FindAvailable(context.Context, *Specification, *Response) error
	// 创建货轮服务
	Create(context.Context, *Vessel, *Response) error
}
*/
import (
	"context"

	pd "github.com/chenxull/shippy/vessel-service/proto/vessel"

	"gopkg.in/mgo.v2"
)

type service struct {
	session *mgo.Session
}

func (s *service) GetRepo() Repository {
	return &VesselRepository{s.session.Clone()}
}

func (s *service) Create(ctx context.Context, req *pd.Vessel, resp *pd.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	if err := repo.Create(req); err != nil {
		return err
	}
	resp.Vessel = req
	resp.Create = true
	return nil
}

func (s *service) FindAvailable(ctx context.Context, req *pd.Specification, resq *pd.Response) error {
	defer s.GetRepo().Close()
	//寻找可用的vessel
	vessel, err := s.GetRepo().FindAvailable(req)
	if err != nil {
		return err
	}
	resq.Vessel = vessel
	return nil
}
