package main

import (
	mcfg "go-micro-learning/UserExamples/gateway_api/config"
	"fmt"
	"github.com/micro/go-micro/config/source/etcd"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
	"github.com/opentracing/opentracing-go"
	"go-micro-learning/UserExamples/basis_lib"
	"go-micro-learning/UserExamples/basis_lib/configuration"
	"go-micro-learning/UserExamples/basis_lib/log"
	tracer "go-micro-learning/UserExamples/basis_lib/tracer/jaeger"
	oti "go-micro-learning/UserExamples/basis_lib/tracer/opentracing"
	"go-micro-learning/UserExamples/gateway_api/plugins/accessLog"
	"os"
	//ph "github.com/afex/hystrix-go/hystrix"
)



var (
	fileCfg  *mcfg.FileCfg
)

func init(){
	initCfg()
	// 链路跟踪
	regTracerErr := plugin.Register(plugin.NewPlugin(
		plugin.WithName("tracer"),
		plugin.WithHandler(
			oti.TracerHttpWrapper,
		),
	))
	if regTracerErr != nil {
		fmt.Println("regErr 链路跟踪 失败")
		os.Exit(1)
	}

	// 新log记录
	regAuthTokenErr := plugin.Register(plugin.NewPlugin(
		plugin.WithName("authTokenJwt"),
		plugin.WithHandler(
			accessLog.FromJWTAuthWrapper(),
		),
	))
	if regAuthTokenErr != nil {
		fmt.Println("regErr 注册JWT验证 失败")
		os.Exit(1)
	}


}

func main() {
	oti.SetSamplingFrequency(100)
	t, io, err := tracer.NewTracer(fileCfg.AppCfg.Name, fileCfg.Trace.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 通过HTTP公开仪表板指标的服务
	//hystrixStreamHandler := ph.NewStreamHandler()
	//hystrixStreamHandler.Start()
	//go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
	cmd.Init()
}






func initCfg(){
	mcfg.InitFileCfg()
	fileCfg = mcfg.GetFileCfg()
	etcdSource := etcd.NewSource(
		etcd.WithAddress(fileCfg.CfgEtcd.Addr...),
		etcd.WithPrefix(fileCfg.CfgEtcd.PathPrefix),
		etcd.StripPrefix(false),
		etcd.WithDialTimeout(2000000),
		)
	basis_lib.Init(
		configuration.WithPathPrefix(fileCfg.CfgEtcd.PathPrefix),
		configuration.WithSource(etcdSource),
		)

	log.Warn("所有加载完成")
}