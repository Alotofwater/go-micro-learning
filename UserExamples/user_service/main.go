package main

import (
	mcfg "go-micro-learning/UserExamples/user_service/config"
	"fmt"
	"go-micro-learning/UserExamples/basis_lib/configuration"
	oti "go-micro-learning/UserExamples/basis_lib/tracer/opentracing"
	"time"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	etcdCfg "github.com/micro/go-micro/config/source/etcd"
	"github.com/micro/go-micro/registry"
	etcdReg "github.com/micro/go-micro/registry/etcd"
	"go-micro-learning/UserExamples/basis_lib/log"
	"go-micro-learning/UserExamples/basis_lib"
	tracer "go-micro-learning/UserExamples/basis_lib/tracer/jaeger"
	"go-micro-learning/UserExamples/user_service/handler"
	"go-micro-learning/UserExamples/user_service/model"
	s "go-micro-learning/UserExamples/user_service/proto/user"
	mdlw "go-micro-learning/UserExamples/user_service/middleware"
	"github.com/opentracing/opentracing-go"
)

var (
	fileCfg  *mcfg.FileCfg
)




func main() {
	// 初始化配置、数据库等信息
	initCfg()

	// 使用 Etcd 注册
	micReg := etcdReg.NewRegistry(registryOptions)
	// 链路跟踪
	t, io, err := tracer.NewTracer(fileCfg.AppCfg.Name, fileCfg.Trace.Addr)
	if err != nil {
		log.Fatal("NewTracer err", err)
	}
	defer io.Close()
	fmt.Println("t, io, err",t, io, err)
	opentracing.SetGlobalTracer(t)
	// 新建服务
	service := micro.NewService(

		micro.Name(fmt.Sprintf("%v.%v",fileCfg.AppCfg.Namespace,fileCfg.AppCfg.Name)),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(micReg),
		micro.Version("latest"),
		micro.WrapHandler(oti.NewGrpcHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(mdlw.AccessLogHandlerWrapper()),
		micro.WrapCall(oti.NewGrpcCallWrapper(opentracing.GlobalTracer())),
	)


	// 服务初始化
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 初始化模型层
			model.Init()
			// 初始化handler
			handler.Init()
		}),
	)

	// 注册服务
	s.RegisterUserHandler(service.Server(), new(handler.Service))

	// 启动服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
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
	log.Warn("所有加载完成")
}