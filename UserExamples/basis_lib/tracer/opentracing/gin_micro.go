package opentracing

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	//"go-micro-learning/UserExamples/basis_lib/log"
	"math/rand"
	"net/http"
)

var (
	contextTracerKey = "Tracer-context"
)

// TracerWrapper tracer 中间件
func TracerGinWrapper(c *gin.Context) {
	//log.Warn()
	md := make(map[string]string)
	headTraceId,headTraceOk = GetGinHeadTraceId(c.Request)
	// 请求体上下文 添加 key value
	if headTraceOk == false {
		headTraceId = GenTraceId()
		//ctx := context.WithValue(context.Background(), TraceIdKey{}, headTraceId)
		//c.Request = c.Request.WithContext(ctx)
	}


	sp := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
	tracer := opentracing.GlobalTracer()

	nsf := sf
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))

	if err == nil { // 如果 为 nil 说明 获取到 父 Tracer 配置
		sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx)) // 创建子配置
		tracer = sp.Tracer()
		nsf = 100
	}
	defer sp.Finish()



	sp.SetTag("Dy-Trace-Id",headTraceId)
	c.Request.Header.Set("Dy-Trace-Id",headTraceId)
	md["Dy-Trace-Id"] = headTraceId



	if err := tracer.Inject(
		sp.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(md)); err != nil {
		//log.Error(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, TraceIdKey{}, headTraceId)  // 创建 TraceId
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)  // 放入 grpc ctx
	// 放入 请求体 ctx
	c.Request = c.Request.WithContext(ctx)

	c.Set(contextTracerKey, ctx)

	c.Next()

	statusCode := c.Writer.Status()
	ext.HTTPStatusCode.Set(sp, uint16(statusCode))
	ext.HTTPMethod.Set(sp, c.Request.Method)
	ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
	if statusCode >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	} else if rand.Intn(100) > nsf {
		ext.SamplingPriority.Set(sp, 0)
	}
}

// ContextWithSpan 返回context
func ContextGinWithSpan(c *gin.Context) (ctx context.Context, ok bool) {
	v, exist := c.Get(contextTracerKey)
	//c.Request.Context().Value()
	if exist == false {
		ok = false
		ctx = context.TODO()
		return
	}

	ctx, ok = v.(context.Context)
	return
}

