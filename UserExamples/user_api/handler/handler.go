package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"

	hystrix_go "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"go-micro-learning/UserExamples/basis_lib/log"
	us "go-micro-learning/UserExamples/user_service/proto/user"
)

var (
	serviceClient us.UserService
)

// Error 错误结构体
type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func Init() {
	// 熔断配置
	hystrix_go.DefaultVolumeThreshold = 1
	hystrix_go.DefaultErrorPercentThreshold = 1
	cl := hystrix.NewClientWrapper()(client.DefaultClient)
	serviceClient = us.NewUserService("go.micro.tc.user-service", cl)
}

// 查询用户
func Login(c *gin.Context) {
	// 链路跟踪
	ctxReq := c.Request.Context()
	//ctx, ok := oti.ContextGinWithSpan(c)
	//if ok == false {
	//	log.Error("gin2micro get context err ",ctx,ok)
	//}

	// 调用后台服务
	masd := &us.Request{
		UserName: c.Query("user_name"),
	}
	rsp, err := serviceClient.QueryUserByName(ctxReq, masd)

	if err != nil {
		log.Error("serviceClient err ", err )
	}
	log.WarnWith(ctxReq,"rsp  ",rsp.User.Name,rsp.Success)
	//if err := c.ShouldBindJSON(&serviceClient); err != nil {
	//	c.AbortWithError(http.StatusBadRequest, errors.New("JWT decode failed"))
	//	return
	//}



	c.JSON(http.StatusCreated, gin.H{"xx":"xxxc查询用户"})
}
