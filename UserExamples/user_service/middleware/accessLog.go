package middleware

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/server"
	"go-micro-learning/UserExamples/basis_lib/log"
	"time"
)

func AccessLogHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			timeStart := time.Now()

			fmt.Println("zxczxc",)
			err := h(ctx, req, rsp)
			if err != nil {
				return err
			}
			timeElapsed := fmt.Sprintf("%v",float64(time.Now().Sub(timeStart).Microseconds()) / 1000)
			log.AcGrpcDebug(ctx, timeElapsed)
			return nil
		}
	}
}
