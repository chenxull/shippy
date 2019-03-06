build:
	#使用micro插件来编译protobuf文件
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/chenxull/shippy/consignment-service proto/consignment/consignment.proto
	#告知GO编译器生成二进制文件的目标环境
	GOOS=linux GOARCH=amd64 go build
	#根据当前目录下的Dockerfile生成名为consignment的镜像
	docker build -t consignment-service .

run:
	# 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
    # 可添加 -d 参数将微服务放到后台运行
	docker run --rm -p 50051:50051 \
	 -e MICRO_SERVER_ADDRESS=:50051 \
	 -e MICRO_REGISTRY=mdns \
	 consignment-service