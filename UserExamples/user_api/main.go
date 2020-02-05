package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/client"
	etcdCfg "github.com/micro/go-micro/config/source/etcd"
	"go-micro-learning/UserExamples/basis_lib/configuration"
	oti "go-micro-learning/UserExamples/basis_lib/tracer/opentracing"
	mcfg "go-micro-learning/UserExamples/user_api/config"
	"go-micro-learning/UserExamples/user_api/routers"
	"os"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/cli"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/web"
	hystrixplugin "github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"github.com/opentracing/opentracing-go"
	"go-micro-learning/UserExamples/basis_lib"
	tracer "go-micro-learning/UserExamples/basis_lib/tracer/jaeger"
	"go-micro-learning/UserExamples/user_api/handler"
)

var (
	fileCfg  *mcfg.FileCfg
)

func main() {
	// 初始化配置
	initCfg()

	// 使用etcd注册
	micReg := etcd.NewRegistry(registryOptions)

	// 链路跟踪
	t, io, err := tracer.NewTracer(fileCfg.AppCfg.Name, fileCfg.Trace.Addr)
	if err != nil {
		fmt.Println("NewTracer err ", err)
		os.Exit(1)
	}
	defer io.Close()
	//设置采样率
	oti.SetSamplingFrequency(100)
	opentracing.SetGlobalTracer(t)


	// 创建新服务
	service := web.NewService(
		web.Name(fmt.Sprintf("%v.%v",fileCfg.AppCfg.Namespace,fileCfg.AppCfg.Name)),
		web.Version(fileCfg.AppCfg.Version),
		web.RegisterTTL(time.Second*15),
		web.RegisterInterval(time.Second*10),
		web.Registry(micReg),
		//web.MicroService(grpc.NewService()),
	)

	// 初始化服务
	if err := service.Init(
		web.Action(
			func(c *cli.Context) {
				// 初始化handler
				handler.Init()
			}),
	); err != nil {
		fmt.Println("service.Init err ", err)
		os.Exit(1)
	}


	// 熔断
	hystrix.DefaultTimeout = 5000
	sClient := hystrixplugin.NewClientWrapper()(service.Options().Service.Client()) // 接口 断言
	errSClient := sClient.Init(
		// client.WrapCall(plugin.NewCallWrapper(t)),
		// 重试次数
		client.Retries(3),
		// 重试 设置重试时要使用的重试函数。
		client.Retry(func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
			fmt.Println(req.Method(), retryCount, " client retry")
			return true, nil
		}),
	)
	if errSClient != nil {
		fmt.Println("errSClient", errSClient)
		os.Exit(1)
	}

	// gin
	router :=routers.InitRouter(fileCfg.AppCfg.Name)
	//urlPrefix := fileCfg.AppCfg.Name
	//router := gin.Default() // 创建默认路由
	//
	//apiv1 := router.Group(fmt.Sprintf("/%v/api/v1",urlPrefix)) //
	//apiv1.Use(oti.TracerGinWrapper) // 链路跟踪
	//
	//apiv1.GET("/user", func(context *gin.Context) {
	//	context.JSON(200,"mxdaa")
	//})
	//// 登录接口
	//apiv1.GET("/user/login",handler.Login)

	////hystrixStreamHandler := hystrix.NewStreamHandler()
	//hystrixStreamHandler.Start()
	//go http.ListenAndServe(net.JoinHostPort("", "33244"), hystrixStreamHandler)

	service.Handle("/",router)  //
	// 运行服务
	if err := service.Run(); err != nil {
		fmt.Println("gin router", errSClient)
		os.Exit(1)
	}
}

func registryOptions(ops *registry.Options) {
	ops.Addrs = fileCfg.RegEtcd.Addr
}

func initCfg(){
	mcfg.InitFileConf()
	fileCfg = mcfg.GetFileCfg()
	etcdSource := etcdCfg.NewSource(
		etcdCfg.WithAddress(fileCfg.CfgEtcd.Addr...),
		etcdCfg.WithPrefix(fileCfg.CfgEtcd.PathPrefix),
		etcdCfg.StripPrefix(false),
		etcdCfg.WithDialTimeout(2000000),
	)
	basis_lib.Init(
		configuration.WithPathPrefix(fileCfg.CfgEtcd.PathPrefix),
		configuration.WithSource(etcdSource),
	)
	fmt.Println("所有加载完成")
}