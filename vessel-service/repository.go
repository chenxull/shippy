package main

import (
	pd "github.com/chenxull/shippy/vessel-service/proto/vessel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName           = "shippy"
	vesselCollection = "vessels"
)

// Repository 货物存储接口
type Repository interface {
	FindAvailable(*pd.Specification) (*pd.Vessel, error)
	Create(*pd.Vessel) error //存放新货物
	Close()
}

// VesselRepository 货物托运服务数据结构，实现了Repository接口
type VesselRepository struct {
	session *mgo.Session
}

// FindAvailable 检查轮船是否可用
func (repo *VesselRepository) FindAvailable(spec *pd.Specification) (*pd.Vessel, error) {
	var vessel *pd.Vessel

	// Here we define a more complex query than our consignment-service's
	// GetAll function. Here we're asking for a vessel who's max weight and
	// capacity are greater than and equal to the given capacity and weight.
	// We're also using the `One` function here as that's all we want.
	err := repo.collection().Find(bson.M{
		"capacity":  bson.M{"$gte": spec.Capacity},
		"maxweight": bson.M{"$gte": spec.MaxWeight},
	}).One(&vessel)
	if err != nil {
		return nil, err
	}
	return vessel, nil

}

// Create 接口实现
func (repo *VesselRepository) Create(c *pd.Vessel) error {
	return repo.collection().Insert(c)
}

// Close 关闭连接
func (repo *VesselRepository) Close() {
	// Close() 会在每次查询结束的时候关闭会话
	// Mgo 会在启动的时候生成一个 "主" 会话
	// 你可以使用 Copy() 直接从主会话复制出新会话来执行，即每个查询都会有自己的数据库会话
	// 同时每个会话都有自己连接到数据库的 socket 及错误处理，这么做既安全又高效
	// 如果只使用一个连接到数据库的主 socket 来执行查询，那很多请求处理都会阻塞
	// Mgo 因此能在不使用锁的情况下完美处理并发请求
	// 不过弊端就是，每次查询结束之后，必须确保数据库会话要手动 Close
	// 否则将建立过多无用的连接，白白浪费数据库资源
	repo.session.Close()
}

// 返回所有货物信息
func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(vesselCollection)
}
