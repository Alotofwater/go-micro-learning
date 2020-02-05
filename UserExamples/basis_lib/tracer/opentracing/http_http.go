package opentracing

// 网关下面 挂 http api 接口  非 grpc接口

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	statusCode "go-micro-learning/UserExamples/basis_lib/breaker/http"
	"math/rand"
	"net/http"
)



// TracerWrapper tracer wrapper
func TracerHttpWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


		headTraceId,headTraceOk = GetGinHeadTraceId(r)
		// 请求体上下文 添加 key value
		if headTraceOk == false {
			headTraceId = GenTraceId()
		}
		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		opts := []opentracing.StartSpanOption{
			opentracing.ChildOf(spanCtx),
			ext.SpanKindRPCClient,
			// opentracing.Tag{Key: string(ext.Component), Value: "koala_rpc"},
			// opentracing.Tag{Key: "ztc", Value: "55555"},
		}
		//sp := opentracing.GlobalTracer().StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
		sp := opentracing.GlobalTracer().StartSpan(r.URL.Path, opts...)
		defer sp.Finish()
		sp.SetTag("Dy-Trace-Id",headTraceId)
		r.Header.Set("Dy-Trace-Id",headTraceId)

		// 请求体上下文添加 key value
		ctx := context.WithValue(context.Background(), TraceIdKey{}, headTraceId)
		r = r.WithContext(ctx)
		// 将 span 信息注入 请求头中
		if err := opentracing.GlobalTracer().Inject(
			sp.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
			//log.Error(" opentracing.GlobalTracer().Injec ",err)
		}

		//opentracing.ContextWithSpan(ctx,sp)

		sct := &statusCode.StatusCodeTracker{ResponseWriter: w, Status: http.StatusOK}

		h.ServeHTTP(sct.WrappedResponseWriter(), r)

		ext.HTTPMethod.Set(sp, r.Method)
		ext.HTTPUrl.Set(sp, r.URL.EscapedPath())
		ext.HTTPStatusCode.Set(sp, uint16(sct.Status))
		if sct.Status >= http.StatusInternalServerError {
			ext.Error.Set(sp, true)
		} else if rand.Intn(100) > sf {
			ext.SamplingPriority.Set(sp, 0)
		}
	})
}
