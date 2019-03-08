package main

import (
	pd "github.com/chenxull/shippy/consignment-service/proto/consignment"
	"gopkg.in/mgo.v2"
)

const (
	dbName                = "shippy"
	consignmentCollection = "consignments"
)

// Repository 货物存储接口
type Repository interface {
	Create(*pd.Consignment) error       //存放新货物
	GetALL() ([]*pd.Consignment, error) //获取仓库中的所有货物
	Close()
}

// ConsignmentRepository 货物托运服务数据结构，实现了Repository接口
type ConsignmentRepository struct {
	session *mgo.Session
}

// Create 接口实现
func (repo *ConsignmentRepository) Create(c *pd.Consignment) error {
	return repo.collection().Insert(c)
}

// GetALL 获取全部数据
func (repo *ConsignmentRepository) GetALL() ([]*pd.Consignment, error) {
	var cons []*pd.Consignment
	// Find() 一般用来执行查询，如果想执行 select * 则直接传入 nil 即可
	// 通过 .All() 将查询结果绑定到 cons 变量上
	// 对应的 .One() 则只取第一行记录
	err := repo.collection().Find(nil).All(&cons)
	return cons, err
}

// Close 关闭连接
func (repo *ConsignmentRepository) Close() {
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
func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(consignmentCollection)
}
