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
	"github.com/dgrijalva/jwt-go"
	"time"
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
		Token struct{
			AccessToken string `json:"access_token"`
			ExpiresAt int64 `json:"expires_at"` // 过期时间
			TimeStamp int64 `json:"time_stamp"` // 时间戳
		}
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
		web.Address("0.0.0.0:8082"),
	)
	router := gin.Default() // 创建默认路由

	apiv1 := router.Group("/api/v1")

	userPasswd := &User{Name:"mxd",Password: "qwe123a"}

	apiv1.GET("/user", func(context *gin.Context) {
		context.JSON(200,"mxdaa")
	})
	apiv1.GET("/user/login", func(context *gin.Context) {
		user := context.Query("user")
		passwd := context.Query("passwd")
		fmt.Println("  passwd  ",user,passwd)
		if passwd != userPasswd.Password || user != userPasswd.Name {
			context.JSON(200,"pass err")
			return
		}

		expired := time.Now().Add( 148 * time.Hour).Unix()

		//token := jwt.New(jwt.SigningMethodHS256) // 生成 token jwt 方式1
		//claim["expired"] = expired // 非标准Claims
		//claim["timestamp"] = time.Now().Unix() // 非标准Claims
		//claim := make(jwt.MapClaims) // 放凭证
		//token.Claims = claim

		claims := 	& jwt.StandardClaims{ // 标准 Claims
			ExpiresAt: time.Now().Add(30 * time.Second).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims) // 生成 token jwt 方式2




		accessToken, err := token.SignedString([]byte("mxd.sign"))
		if err != nil {
			context.JSON(200,"accessToken err")

			return
		}
		userPasswd.Token.ExpiresAt = expired
		userPasswd.Token.AccessToken = accessToken
		userPasswd.Token.TimeStamp = time.Now().Unix()

		context.JSON(200,fmt.Sprintf("user:%v  token:%v",user,userPasswd))
	})

	apiv2 := router.Group("/api/v1")
	apiv2.Use(TokenValid)
	apiv2.POST("/user/tokenvalid", func(context *gin.Context) {
		context.JSON(200,fmt.Sprintf("通过"))
		return
	})

	service.Handle("/",router)  //

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// jwt 验证

func TokenValid(context *gin.Context){
	authorization := context.GetHeader("Authorization")
	token, err := jwt.Parse(authorization, func(token *jwt.Token) (i interface{}, e error) {
		return []byte("mxd.sign"),nil
	})
	if err != nil {
		fmt.Println("zxczxc")
		if  err, ok := err.(*jwt.ValidationError); ok  { //  强转 jwt 错误类型
			if err.Errors & jwt.ValidationErrorMalformed != 0 {
				// 验证不通过，不再调用后续的函数处理
				context.Abort()
				context.JSON(200,fmt.Sprintf("ValidationErrorMalformed cccc %v",err))
				return
			}
			if err.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
				// 验证不通过，不再调用后续的函数处理
				context.Abort()
				context.JSON(200,fmt.Sprintf("ValidationErrorExpired xxx ValidationErrorNotValidYet%v",err))

				return
			}
			context.JSON(200,fmt.Sprintf("wi  %v",err))
			return
		}
	}
	if token.Valid {
		//context.Next()
		return
	}else {
		// 验证不通过，不再调用后续的函数处理
		context.Abort()

		context.JSON(200,fmt.Sprintf("no token.Valid  %v"))
	}
}




