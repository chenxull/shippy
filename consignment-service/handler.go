package main

import (
	"context"
	"fmt"
	"log"

	pd "github.com/chenxull/shippy/consignment-service/proto/consignment"
	vesselProto "github.com/chenxull/shippy/vessel-service/proto/vessel"
	"gopkg.in/mgo.v2"
)

type handler struct {
	session      *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

// GetRepo 从主回话中clone()出新的回话处理查询
func (h *handler) GetRepo() Repository {
	return &ConsignmentRepository{h.session.Clone()}
}

//托运新的货物
func (h *handler) CreateConsignment(ctx context.Context, req *pd.Consignment, resp *pd.Response) error {
	repo := h.GetRepo()
	defer repo.Close()
	// 检查是否有合适的货轮
	// 通过创建一个货轮服务的client实例来发出请求查询是否有合适的货轮
	vReq := &vesselProto.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vesselResp, err := h.vesselClient.FindAvailable(context.Background(), vReq)
	log.Printf("Found vessel: %s \n", vesselResp.Vessel.Name)
	if err != nil {
		return err
	}

	// 将从货轮服务中获取的货轮号赋值给req
	req.VesselId = vesselResp.Vessel.Id
	//存储承运的货物
	err = repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	resp.Created = true
	resp.Consignment = req
	fmt.Println(resp)
	return nil
}

// 获取目前所托运的货物
func (h *handler) GetConsignment(ctx context.Context, req *pd.GetRequest, resp *pd.Response) error {
	repo := h.GetRepo()
	defer repo.Close()
	consignments, err := repo.GetALL()
	if err != nil {
		return err
	}
	//resp = &pd.Response{Consignments: allconsignments}
	resp.Consignments = consignments
	return nil
}
