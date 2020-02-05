package main

import (
	"context"
	"fmt"
	"math/rand"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-micro/registry"
	hello "go-micro-learning/protoc"
	"log"
	"strconv"
)

type Say struct{
	Tag string
}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name + "[ From " + s.Tag + "]"
	return nil
}

func main() {
	etcdRegistry := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = []string{"127.0.0.1:2379"}
			//etcdv3.Auth("sss","xxx")() // 密码
		})

	metdataMap := map[string]string{"rmb":"9999"}


	// 创建服务
	service := micro.NewService(
		micro.Name("go.micro.api.greeter"),
		micro.Registry(etcdRegistry),
		micro.Version("mxd00010"), // 修改版本信息
		micro.Metadata(metdataMap), // 修改当前服务 metadata
		micro.Transport(grpc.NewTransport()), // 当前服务传输协议 (与客户端匹配)

		//micro.RegisterTTL(time.Second*30),
		//micro.RegisterInterval(time.Second*10),
	)
	// 负载均衡

	// optionally setup command line usage
	service.Init()

	// Register Handlers

	say:= &Say{
		Tag:strconv.Itoa(rand.Int()), // 随机数
	}
	fmt.Println("当前服务Tag为 ", say.Tag)
	hello.RegisterSayHandler(service.Server(), say)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}