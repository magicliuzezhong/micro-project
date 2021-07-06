//
// Package web_register
// @Description：web初始化
// @Author：liuzezhong 2021/6/25 6:16 下午
// @Company cloud-ark.com
//
package web_register

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"micro-project/internal/pkg/filter"
	"micro-project/internal/pkg/jeager"
	"micro-project/internal/pkg/util"
)

//
// InitWeb
// @Description: 初始化web服务
// @param exitCallback 退出回调
//
func InitWeb(exitCallback func(), extend func(application *iris.Application)) {
	app := iris.New()
	app.Use(customRecover)            //注册服务恢复逻辑，全局数据类型统一格式返回
	app.Use(filter.LimiterFilter)     //注册限流器
	app.Use(filter.WebFilter)         //注册web拦截器
	app.Use(jeager.JeagerTraceFilter) //注册jeager全链路追踪器

	app.Handle("GET", "/health", func(context context.Context) {
		_, _ = context.ResponseWriter().Write([]byte(`{"status": "ok"}`))
	})

	app.Handle("GET", "/metrics", iris.FromStd(promhttp.Handler())) //集成prometheus

	if extend != nil {
		extend(app)
	}

	iris.RegisterOnInterrupt(func() { //监听退出事件，程序中断执行此处
		if exitCallback != nil {
			exitCallback()
		}
	})

	var err = app.Run(iris.Addr(":"+util.GetApplication().Server.Port),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithCharset("UTF-8"),
	)
	if err != nil {
		panic("====微服务启动异常====")
	}
}
