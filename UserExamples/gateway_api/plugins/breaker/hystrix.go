package breaker

import (
	"errors"
	"fmt"
	"go-micro-learning/UserExamples/basis_lib/log"
	"net/http"
	statusCode "go-micro-learning/UserExamples/basis_lib/breaker/http"
	"github.com/afex/hystrix-go/hystrix"
)

// BreakerWrapper hystrix breaker
func BreakerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + "-" + r.RequestURI
		//log.Warn("BreakerWrapper: ",name)
		hystrix.DefaultVolumeThreshold = 300
		hystrix.DefaultTimeout = 10000
		hystrix.DefaultMaxConcurrent = 2000
		hystrix.DefaultErrorPercentThreshold = 0
		err := hystrix.Do(name, func() error {
			sct := &statusCode.StatusCodeTracker{ResponseWriter: w, Status: http.StatusOK}
			h.ServeHTTP(sct.WrappedResponseWriter(), r)

			if sct.Status >= http.StatusBadRequest {
				str := fmt.Sprintf("status code %d", sct.Status)
				log.Warn(str)
				return errors.New(str)
			}
			return nil
		}, nil)
		if err != nil {
			log.Error("BreakerWrapper: hystrix breaker err: ", err)
			return
		}
	})
}


