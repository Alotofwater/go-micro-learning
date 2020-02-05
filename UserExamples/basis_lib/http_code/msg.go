package mcode

import (
	"github.com/gin-gonic/gin"
	"go-micro-learning/UserExamples/basis_lib/log"
	"net/http"
)

var msgFlags = map[int]string {
	SUCCESS : "ok",
	ERROR : "fail",
	INVALID_PARAMS : "missing parameters | 请求参数错误或参数不存在",
	METHOD_NOT_ALLOWED : "method not allowed | 方法不允许",
	NO_ROUTE : "No Route| 路由未找到",
	ERROR_EXIST_TAG : "已存在该标签名称",
	ERROR_NOT_EXIST_TAG : "该标签不存在",
	ERROR_NOT_EXIST_APP_PYPE : "不存在的应用类型",
	ERROR_AUTH_CHECK_TOKEN_FAIL : "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT : "Token已超时",
	ERROR_AUTH_TOKEN : "Token生成失败",
	ERROR_AUTH : "Token错误",
	ERROR_CELERY_TIOMEOUT: "Celery 请求超时",
	ERROR_NOT_FOUND: "访问不存在",
	SERVICE_FRAMEWORK_ERROR: "服务框架错误",
}


func GetResponseMsg(code int, data interface{},c... *gin.Context) (map[string]interface{}) {
	//fmt.Println("c...",c[0].Request.URL.Path,len(c))
	responseMap := make(map[string]interface{})
	msg, ok := msgFlags[code]
	if ok {
		responseMap["code"] = code
		responseMap["msg"] = msg
		responseMap["data"] = data
		//responseMap["request"] = c[0].Request.URL.Path
		return responseMap
	}

	responseMap["code"] = code
	responseMap["msg"] = "MsgFlags中无该错误码！"
	responseMap["data"] = data
	//responseMap["request"] = c[0].Request.URL.Path

	return responseMap
}


func errlog(requestURL string,resp interface{}){
	mapVal,ok := resp.( map[string]interface{})
	if ok {
		log.Errorf("errlogResp:%+v   requestURL:%+v",mapVal["error"],requestURL)
		delete(mapVal,"error")
	}
}
// 参数不存在
func InvalidParams(c *gin.Context, resp interface{}) {
	//c.JSON(http.StatusBadRequest, gin.H{"code": "4000", "msg": msgFlags["4000"], "data": resp, "request": c.Request.URL.Path})
	c.JSON(http.StatusOK, gin.H{"code": "4000", "msg": msgFlags[4000], "data": resp, "request": c.Request.URL.Path})
	return
}

// 成功响应
func Success(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": "2000", "msg": msgFlags[2000], "data": resp, "request": c.Request.URL.Path})
	//c.JSON(http.StatusOK, gin.H{"errcode": 0, "resp": resp, "request": c.Request.URL.Path})
	return
}


// 没有找到
func NotFound(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": "42000", "msg": msgFlags[42000], "data": resp, "request": c.Request.URL.Path})
	return
}

// 服务错误
func ServerError(c *gin.Context, resp interface{}) {
	requestURL := c.Request.URL.Path
	errlog(requestURL,resp)
	//c.JSON(http.StatusInternalServerError, gin.H{"code": "5000", "msg": msgFlags["5000"], "data": resp, "request": requestURL})
	c.JSON(http.StatusOK, gin.H{"code": "5000", "msg": msgFlags[5000], "data": resp, "request": requestURL})
	return
}



// 服务框架错误 ( panic 致命错误 导致业务逻辑 不能正常执行)
func ServiceFrameworkError(c *gin.Context, resp interface{}) {
	requestURL := c.Request.URL.Path
	errlog(requestURL,resp)
	c.JSON(http.StatusOK, gin.H{"code": "5555", "msg": msgFlags[5555], "data": resp, "request": requestURL})
	return
}


// 方法不允许
func MethodNotAllowed(c *gin.Context, resp interface{}) {
	c.JSON(403, gin.H{"code": "403", "msg": msgFlags[403], "data": resp, "request": c.Request.URL.Path})
	return
}

// 未找到路由
func NoRoute(c *gin.Context, resp interface{}) {
	c.JSON(404, gin.H{"code": "404", "msg": msgFlags[404], "data": resp, "request": c.Request.URL.Path})
	return
}

