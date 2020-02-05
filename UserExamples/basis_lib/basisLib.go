package basis_lib

import (
	"fmt"
	"go-micro-learning/UserExamples/basis_lib/configuration"
	"go-micro-learning/UserExamples/basis_lib/log"
)

var (
	pluginFuncs []func()
)


func Init(opts ...configuration.Option) {
	// 初始化配置
	fmt.Println("初始化配置 Init")
	configuration.Init(opts...)
	log.InitLog()
	// 加载依赖配置的插件
	for _, f := range pluginFuncs {
		f()
	}
}

func Register(f func()) {
	pluginFuncs = append(pluginFuncs, f)
}

