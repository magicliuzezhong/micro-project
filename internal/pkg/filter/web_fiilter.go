//
// Package filter
// @Description：web过滤器
// @Author：liuzezhong 2021/6/25 6:01 下午
// @Company cloud-ark.com
//
package filter

import (
	"encoding/json"
	"github.com/kataras/iris/v12/context"
	"golang.org/x/time/rate"
	"micro-project/internal/pkg/common"
)

var rateLimiter = rate.NewLimiter(1, 5)

//
// LimiterFilter
// @Description: 限流过滤器
// @param context iris上下文
//
func LimiterFilter(context context.Context) {
	if !rateLimiter.Allow() {
		var response = context.ResponseWriter()
		var resultStr, _ = json.Marshal(common.ResponseResult{
			Status: 429,
			Msg:    "请求过多，请稍后重试",
			Data:   "",
		})
		response.Write(resultStr)
		response.WriteHeader(429)
		context.StopExecution()
		return
	}
	context.Next()
}

//
// WebFilter
// @Description: 普通web过滤器
// @param context iris上下文
//
func WebFilter(context context.Context) {
	var url = context.Request().URL
	if url.String() == "/health" || url.String() == "/metrics" { //如果是健康检查或者普罗米修斯那么不进行全链路追踪
		context.Next()
		return
	}
	context.Next()
}
