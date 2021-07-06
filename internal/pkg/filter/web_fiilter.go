//
// Package filter
// @Description：web过滤器
// @Author：liuzezhong 2021/6/25 6:01 下午
// @Company cloud-ark.com
//
package filter

import (
	"github.com/kataras/iris/v12/context"
	"golang.org/x/time/rate"
	"micro-project/internal/pkg/common"
)

//
// @Description: 限流器（令牌桶算法）
// @param 5000 每秒添加5000个令牌
// @param 20000 总计令牌数
//
var rateLimiter = rate.NewLimiter(1, 5)

//
// LimiterFilter
// @Description: 限流过滤器
// @param context iris上下文
//
func LimiterFilter(context context.Context) {
	if !rateLimiter.Allow() {
		context.Values().Set("val", common.ResponseResult{
			Status: 429,
			Msg:    "请求过多，请稍后重试",
			Data:   "",
		})
		context.StatusCode(429)
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
