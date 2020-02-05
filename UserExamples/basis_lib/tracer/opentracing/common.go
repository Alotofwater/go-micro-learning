package opentracing

import (
	"math/rand"
	"time"
)


var (
    sf = 100
	headTraceId string
	headTraceOk bool
)
func init() {
	rand.Seed(time.Now().Unix())
}

// SetSamplingFrequency 设置采样频率
// 0 <= n <= 100
func SetSamplingFrequency(n int) {
	sf = n
}