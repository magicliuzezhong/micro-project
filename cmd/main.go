//
// Package main
// @Description：微服务项目主入口
// @Author：liuzezhong 2021/6/25 5:38 下午
// @Company cloud-ark.com
//
package main

import (
	"github.com/kataras/iris/v12"
	"micro-project/internal/app/config"
	"micro-project/internal/pkg/discover"
	"micro-project/internal/pkg/util"
	"micro-project/internal/pkg/web_register"
)

func main() {
	util.GetLogger().Println("====微服务程序启动====")
	//微服务注册
	discover.Register()

	//iris web服务注册
	web_register.InitWeb(func() {
		//程序结束回调函数，此处将执行注销服务注册服务
		discover.DeRegister()
	}, func(application *iris.Application) {
		//在此处注册mvc，mvc注册逻辑不应该放在pkg包中，因为这个是属于业务逻辑的注册，不是通用逻辑的注册，目前是放在config下
		config.RegisterIrisMvc(application)
		//还可以进行本地过滤器的添加, 如需实现本地过滤器请参考iris官方文档
		// application.Use(filter.WebFilter) //注册web拦截器
	})
}
