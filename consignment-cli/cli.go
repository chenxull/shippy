package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"

	pd "github.com/chenxull/shippy/consignment-service/proto/consignment"
)

const (
	ADDRESS           = "localhost:50051"
	DEFAULT_INFO_FILE = "consignment.json"
)

//读取consignment.json中记录的货物信息
func parseFile(fileName string) (*pd.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pd.Consignment
	//解析data中的json数据，并将结果存储在consignment中，consignment必须是指针类型
	err = json.Unmarshal(data, &consignment)
	//fmt.Println(data)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main() {
	//连接到gRPC服务器
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect error ", err)
	}
	defer conn.Close()

	//初始化gRPC服务器
	client := pd.NewShippingServiceClient(conn)

	//在命令行中指定新的货物信息json文件
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}
	fmt.Println(infoFile)
	//解析货物信息
	consignment, err := parseFile(infoFile)
	if err != nil {
		log.Fatal("parse info file error ", err)
	}

	//使用CreateConsignment创建了一个发送货物的请求
	resp, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatal("create consignment error ", err)
	}
	log.Printf("created:%t", resp.Created)

	// 列出目前所有托运的货物
	resp, err = client.GetConsignment(context.Background(), &pd.GetRequest{})
	if err != nil {
		log.Fatal("failed to list consignment", err)
	}

	for _, c := range resp.Consignments {
		log.Printf("%+v", c)
	}
}
