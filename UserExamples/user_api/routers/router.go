package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	resp_code "go-micro-learning/UserExamples/basis_lib/http_code"
	"go-micro-learning/UserExamples/basis_lib/log"
	oti "go-micro-learning/UserExamples/basis_lib/tracer/opentracing"
	"go-micro-learning/UserExamples/user_api/handler"
	"time"
)


func RecoverErr() gin.HandlerFunc {
	// 接收 框架 报错
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				//mlogrus.LogError("发生错误:%v  方法:%v ", err, c.HandlerName())
				//c.AbortWithStatusJSON(500, gin.H{"error": "sorry, we made a mistake!", "errcode": 1099, "url": c.Request.URL.Path})
				//c.AbortWithStatusJSON(500, gin.H{"error": "sorry, we made a mistake!", "errcode": 1099, "url": c.Request.URL.Path})
				resp_code.ServiceFrameworkError(c,map[string]interface{}{"errorTo":"服务框架错误"})

			}
		}()
		c.Next()
	}
}

func InitRouter(appName string) *gin.Engine {
	log.Warn("zxc")
	r := gin.New() // New返回一个新的未附加任何中间件的空白引擎实例。
	r.Use(oti.TracerGinWrapper) // 链路跟踪
	r.Use(func(c *gin.Context) {
		timeStart := time.Now()
		c.Next()
		timeElapsed := fmt.Sprintf("%v",float64(time.Now().Sub(timeStart).Microseconds()) / 1000)
		log.AcGinDebug(c,timeElapsed)

	})

	r.Use(RecoverErr())
	//r.Use(gin.Logger())
	//r.Use(gin.Recovery())
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		//c.JSON(403, gin.H{"error": "method not allowed", "errcode": 1010, "request": c.Request.URL.Path})
		resp_code.MethodNotAllowed(c, map[string]interface{}{"errorTo": "方法不能执行"})
	})
	r.NoRoute(func(c *gin.Context) {
		//c.JSON(404, gin.H{"error": "route not found", "errcode": 1010, "request": c.Request.URL.Path})
		resp_code.NoRoute(c, map[string]interface{}{"errorTo": "路由未找到"})
	})

	//r.Use(util.MtestMiddleware())

	gin.SetMode("debug") // 运行模式 debug 或者release
	urlPrefix := appName
	apiv1 := r.Group(fmt.Sprintf("/%v/api/v1",urlPrefix)) //
	{
		apiv1.GET("/user", func(context *gin.Context) {
			context.JSON(200,"mxdaa")
		})
		// 登录接口
		apiv1.GET("/user/login",handler.Login)
	}

	return r
}