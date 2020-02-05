package opentracing

import (
	"context"
	"google.golang.org/grpc/metadata"
)
// 中间件
func cc (ctx context.Context){
	//处理traceId
	var traceId string
	//从ctx获取grpc的metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		vals, ok := md["util.TraceID"]
		if ok && len(vals) > 0 {
			traceId = vals[0]
		}
	}

	if len(traceId) == 0 {
		traceId = GenTraceId()
	}

	//ctx = log.WithFieldContext(ctx)
	ctx = WithTraceId(ctx, traceId)
	//resp, err = next(ctx, req)
	return
}
