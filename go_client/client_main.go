package main

import (
	"context"
	"fmt"
	hystrixT "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
	proto "go-micro-learning/protoc"
	"time"
)



type NYclientWrapper struct {
	client.Client
}

func (c *NYclientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrixT.Do(req.Service()+"."+req.Endpoint(), func() error {
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(e error) error {
		fmt.Println("这是一个备用服务器")
		return nil
	})
}

// NewClientWrapper returns a hystrix client Wrapper.
func NYNewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &NYclientWrapper{c}
	}
}

func main(){
	etcdRegistry := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = []string{"127.0.0.1:2379"}
			//etcdv3.Auth("sss","xxx")(options) // 密码
		})
	// 创建一个新的服务
	//hystrixT.DefaultTimeout = 1900 // 熔断超时时间
	service := micro.NewService(
		micro.Name("greeter.clients"),
		micro.Registry(etcdRegistry),
		micro.Transport(grpc.NewTransport()), // 当前服务传输协议 (与 服务端匹配)
		//micro.WrapClient(hystrix.NewClientWrapper()), // "github.com/micro/go-plugins/wrapper/breaker/hystrix"
		micro.WrapClient(NYNewClientWrapper()),
	)
	service.Init()

	// 创建 greeter 客户端
	greeter := proto.NewSayService("go.micro.api.greeter",service.Client()) // go.micro.api.greeter 服务端名称

	t := time.NewTicker(time.Millisecond * 600)
	for e := range t.C {
		rsp, err :=greeter.Hello(context.TODO(),&proto.Request{Name:"xiao zhang"}, func(options *client.CallOptions) {
			options.RequestTimeout = time.Millisecond * 1400 // 请求超时
			// 重试
			//options.Retry = func(ctx context.Context, req client.Request, retryCount int, err error) (b bool, e error) {
			//	fmt.Println("Retry")
			//	return  false, nil
			//}
		})
		if err !=nil {
			fmt.Println("err", err.Error())
			//return
		}else {
			fmt.Printf("rsp.Msg: %v e: %v   \n\n",rsp.Msg,e)

		}
		//fmt.Println("rsp.Msg",rsp.GetMsg())
	}

}
