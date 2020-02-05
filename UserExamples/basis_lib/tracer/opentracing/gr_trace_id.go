package opentracing


import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

const (
	MaxTraceId = 100000000
)

type TraceIdKey struct{}
var (
	HeadTraceIdKey string = "Dy-Trace-Id"
)
func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetCtxTraceId(ctx context.Context) (traceId string) {

	traceId, ok := ctx.Value(TraceIdKey{}).(string)
	if !ok {
		traceId = GenTraceId()
	}

	return
}

func GetGinHeadTraceId(r *http.Request) (traceId string,ok bool) {
	/*
	ok true：说明请求头中有TraceId false：说明请求头中无TraceId
	*/
	traceId = r.Header.Get(HeadTraceIdKey)
	// 请求头 无 TraceId
	if len(traceId) <= 0 {
		ok = false
		traceId = GenTraceId()
		return traceId, ok
	}
	ok = true
	return traceId, ok
}

func GetGrpcHeadTraceId(md map[string]string) (traceId string,ok bool) {
	/*
		ok true：说明请求头中有TraceId false：说明请求头中无TraceId
	*/
	traceId = md[HeadTraceIdKey]
	// 请求头 无 TraceId
	if len(traceId) <= 0 {
		ok = false
		traceId = GenTraceId()
		return traceId, ok
	}
	ok = true
	return traceId, ok
}

func GetHttpHeadTraceId(r *http.Request) (traceId string,ok bool) {
	/*
		ok true：说明请求头中有TraceId false：说明请求头中无TraceId
	*/

	traceId = r.Context().Value(TraceIdKey{}).(string)
	// 上下文 是否 有 TraceIdKey
	if len(traceId) <= 0 {
		ok = false
		traceId = GenTraceId()
		return traceId, ok
	}
	ok = true
	return traceId, ok
}
func GetGinHeadTraceIdvvv(c *gin.Context) (traceId string,ok bool) {
	/*
		ok true：说明请求头中有TraceId false：说明请求头中无TraceId
	*/
	//traceId = c
	//// 请求头 无 TraceId
	//if len(traceId) <= 0 {
	//	ok = false
	//	traceId = GenTraceId()
	//	return traceId, ok
	//}
	//ok = true
	//return traceId, ok
	return
}


func GenTraceId() (traceId string) {
	now := time.Now()
	traceId = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%08d", now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), rand.Int31n(MaxTraceId))
	return
}

func WithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, TraceIdKey{}, traceId)
}


