package main

import (
	"flag"
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	"github.com/gin-gonic/gin"
)



type (
	Config struct {
		Version string
		Hello struct{
			Name string
		}
		Etcd struct{
			Addr []string
			User string
			Passwd string
		}
	}

	User struct {
		Name string `json:"name"`
		Password string `json:"password"`
		
	}

)


func main() {
	configFile := flag.String("f","./config/config.yaml","please use config.yaml")
	conf := &Config{}
	if err := config.LoadFile(*configFile); err != nil {
		fmt.Println("config.LoadFile err ", err)
	}
	if err := config.Scan(conf);err != nil {
		fmt.Println("config.Scan err ", err)

	}
	etcdRegistry := etcdv3.NewRegistry(
		func(options *registry.Options) {
			//options.Addrs = []string{"127.0.0.1:2379"}
			options.Addrs = conf.Etcd.Addr
			//etcdv3.Auth("sss","xxx")() // 密码
		})

	metdataMap := map[string]string{"rmb":"9999"}
	//ratelimit.

	// 创建 web api  服务
	service := web.NewService(
		web.Name("gin.api.server"),
		web.Registry(etcdRegistry),
		web.Version("mxd00010"), // 修改版本信息
		web.Metadata(metdataMap), // 修改当前服务 metadata
		web.Address("0.0.0.0:8080"),
	)
	router := gin.Default() // 创建默认路由
	apiv1 := router.Group("/api/v1")
	{
		apiv1.GET("/user", func(context *gin.Context) {
			context.JSON(200,"mxdaa")
		})
	}
	service.Handle("/",router)  //

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}




