package main

import (
	"fmt"
	"github.com/micro/go-micro/config/source/etcd"
	"github.com/micro/go-plugins/micro/cors"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
	"github.com/opentracing/opentracing-go"
	"go-micro-learning/UserExamples/basis_lib"
	"go-micro-learning/UserExamples/basis_lib/configuration"
	"go-micro-learning/UserExamples/basis_lib/log"
	tracer "go-micro-learning/UserExamples/basis_lib/tracer/jaeger"
	oti "go-micro-learning/UserExamples/basis_lib/tracer/opentracing"
	mcfg "go-micro-learning/UserExamples/gateway_api/config"
	"go-micro-learning/UserExamples/gateway_api/plugins/accessLog"
	"go-micro-learning/UserExamples/gateway_api/plugins/monitoring"
	"os"
	//ph "github.com/afex/hystrix-go/hystrix"
)



var (
	fileCfg  *mcfg.FileCfg
)

func init(){
	initCfg()
	// 链路跟踪
	if err := plugin.Register(plugin.NewPlugin(
		plugin.WithName("tracer"),
		plugin.WithHandler(
			oti.TracerHttpWrapper,
		),
	));err != nil {
		fmt.Println("regErr 链路跟踪 失败 ",err)
		os.Exit(1)
	}
	// 熔断器
	//if err := plugin.Register(plugin.NewPlugin(
	//	plugin.WithName("breaker"),
	//	plugin.WithHandler(
	//		breaker.BreakerWrapper,
	//	),
	//)); err != nil {
	//	fmt.Println("regErr 熔断器 失败 ",err)
	//	os.Exit(1)
	//}
	// 新log记录
	if err := plugin.Register(plugin.NewPlugin(
		plugin.WithName("authTokenJwt"),
		plugin.WithHandler(
			accessLog.AccessLogWrapper(),
		),
	)); err != nil {
		fmt.Println("regErr 注册JWT验证 失败 ",err)
		os.Exit(1)
	}


	//  注册跨域插件
	if err := plugin.Register(cors.NewPlugin()); err != nil {
		fmt.Println("regErr 注册跨域插件 失败 ",err)
		os.Exit(1)
	}

	// 参数形式 prometheus
	//if err := plugin.Register(metrics.NewPlugin()); err != nil {
	//	fmt.Println("regErr 参数形式 prometheus 失败 ",err)
	//	os.Exit(1)
	//}

	if err := plugin.Register(plugin.NewPlugin(
		plugin.WithName("HandlerMetrics"),
		plugin.WithHandler(
			monitoring.MetricsWrapper,
		),
	));  err != nil {
		fmt.Println("regErr  prometheus 失败 ",err)
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