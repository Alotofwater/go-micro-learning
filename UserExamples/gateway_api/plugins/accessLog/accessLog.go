package accessLog

import (
	"fmt"
	"github.com/micro/micro/plugin"
	statusCode "go-micro-learning/UserExamples/basis_lib/breaker/http"
	"go-micro-learning/UserExamples/basis_lib/log"
	"go-micro-learning/UserExamples/basis_lib/sql_db"
	"net/http"
	"time"
)

// 日志
type traceIdKey struct{}
func FromJWTAuthWrapper()plugin.Handler{
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timeStart := time.Now()
			db := sql_db.GetDB()
			row,err := db.Query("select * from rbac_menu")
			row.Close()
			fmt.Println("FromJWTAuthWrapper",row,err)

			// 状态码跟踪
			sct := &statusCode.StatusCodeTracker{ResponseWriter: w, Status: http.StatusOK}
			h.ServeHTTP(sct.WrappedResponseWriter(), r)
			status := sct.Status
			timeElapsed := fmt.Sprintf("%v",float64(time.Now().Sub(timeStart).Microseconds()) / 1000)
			log.AcHttpDebug(r,timeElapsed,status)
			return
		})

	}
}
