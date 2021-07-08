//
// Package web_register
// @Description：iris其他通用处理逻辑
// @Author：liuzezhong 2021/7/6 2:17 下午
// @Company cloud-ark.com
//
package web_register

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
	"micro-project/internal/pkg/common"
	"runtime"
)

//
// checkFilterWhite
// @Description: 检查白名单
// @param ctx 上下文
// @return bool true：是白名单，false：不是白名单值
//
func checkFilterWhite(ctx context.Context) bool {
	var url = ctx.Request().URL
	if url.String() == "/health" || url.String() == "/metrics" {
		ctx.Next()
		return true
	}
	return false
}

//
// buildResult
// @Description: 构建返回结果
// @param code 结果码
// @param msg 信息
// @param data 数据
// @return common.ResponseResult 结果
//
func buildResult(code int, msg string, data interface{}) common.ResponseResult {
	return common.ResponseResult{
		Status: code,
		Msg:    msg,
		Data:   data,
	}
}

//
// customRecover
// @Description: 错误信息处理
// @param ctx 上下文
//
func customRecover(ctx context.Context) {
	if checkFilterWhite(ctx) {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			if ctx.IsStopped() {
				return
			}
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}
				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}
			errMsg := fmt.Sprintf("错误信息: %s", err)
			// when stack finishes
			logMessage := fmt.Sprintf("从错误中恢复：('%s')\n", ctx.HandlerName())
			logMessage += errMsg + "\n"
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			// 打印错误日志
			ctx.Application().Logger().Warn(logMessage)
			ctx.StatusCode(500) // 内部服务错误
			_, _ = ctx.JSON(buildResult(500, fmt.Sprintf("%s", err), ""))
			ctx.StopExecution()
		} else { //此处将统一处理controller的返回数据，如果已经是ResponseResult那么将以传入的为准
			var val = ctx.Values().Get("val")
			var result = buildResult(200, "请求成功", "")
			switch val.(type) {
			case common.ResponseResult: //如果原来就已经是这个结果了那么不处理直接返回
				result = val.(common.ResponseResult)
				break
			case nil:
				break
			default:
				result.Data = val
			}
			_, _ = ctx.JSON(result)
		}
	}()
	ctx.Next()
}
