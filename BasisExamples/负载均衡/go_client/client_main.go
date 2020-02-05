package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-micro/registry"
	proto  "go-micro-learning/protoc"
	"github.com/micro/go-micro"
	"time"
)

func main(){
	etcdRegistry := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = []string{"127.0.0.1:2379"}
			//etcdv3.Auth("sss","xxx")(options) // 密码
		})
	// 创建一个新的服务
	service := micro.NewService(
		micro.Name("greeter.clients"),
		micro.Registry(etcdRegistry),
		micro.Transport(grpc.NewTransport()), // 当前服务传输协议 (与 服务端匹配)

	)
	service.Init()

	// 创建 greeter 客户端
	greeter := proto.NewSayService("go.micro.api.greeter",service.Client())

	t := time.NewTicker(time.Second * 1)
	for e := range t.C {
		rsp, err :=greeter.Hello(context.TODO(),&proto.Request{Name:"xiao zhang"})
		if err !=nil {
			fmt.Println("err", err.Error())
			return
		}
		fmt.Printf("rsp.Msg: %+v e: %+v   \n\n",e,rsp.Msg)
		//fmt.Println("rsp.Msg",rsp.GetMsg())
	}

}
