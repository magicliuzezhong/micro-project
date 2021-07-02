//
// Package config
// @Description：注册iris web mvc
// @Author：liuzezhong 2021/6/28 2:40 下午
// @Company cloud-ark.com
//
package config

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"micro-project/internal/app/controller"
)

//
// RegisterIrisMvc
// @Description: 此处专门用于注册mvc(请勿在此写其他逻辑)
// @param app app
//
func RegisterIrisMvc(app *iris.Application) {
	mvc.New(app.Party("/test")).Handle(new(controller.TestController))
}
