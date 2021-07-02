//
// Package discover
// @Description：服务发现工具类
// @Author：liuzezhong 2021/6/28 11:41 上午
// @Company cloud-ark.com
//
package discover

import (
	"micro-project/internal/pkg/common"
	"micro-project/internal/pkg/util"
	"strconv"
	"sync"
)

//
// @Description: 服务注册客户端
//
var discoverClient IDiscoveryClient

//
// @Description: 用于获取客户端单例
//
var getDiscoverOnce sync.Once

//
// @Description: 用于仅仅注册服务一次
//
var registerOnce sync.Once

//
// getDiscoverClient
// @Description: 获取服务注册实例，该功能目前只实现了consul的注册功能
// @return IDiscoveryClient 服务注册接口
//
func getDiscoverClient() IDiscoveryClient {
	getDiscoverOnce.Do(func() {
		var application = util.GetApplication()
		var ip = application.ConsulConf.Ip
		var port, _ = strconv.Atoi(application.ConsulConf.Port)
		discoverClient = NewConsulDiscoverClient(ip, port)
	})
	return discoverClient
}

//
// Register
// @Description: 服务注册主入口
//
func Register() {
	registerOnce.Do(func() {
		var application = util.GetApplication()
		var ip = application.ConsulConf.Ip
		//var port, _ = strconv.Atoi(application.ConsulConf.Port)
		var port, _ = strconv.Atoi(application.Server.Port)
		var tags = []string{application.ConsulConf.Tag}
		var serviceName = application.Server.ServerName
		var serviceWeight, _ = strconv.Atoi(application.ConsulConf.Weight)
		var registerStatus = getDiscoverClient().Register(ip, port, serviceName, serviceWeight, nil, tags)
		if !registerStatus {
			panic("程序注册失败")
		}
	})
}

//
// DiscoverServices
// @Description: 服务发现
// @param serviceName 服务名称
// @return []*common.ServiceInstance 服务实例
//
func DiscoverServices(serviceName string) []*common.ServiceInstance {
	return getDiscoverClient().DiscoverServices(serviceName)
}

//
// DeRegister
// @Description: 注销服务注册功能
//
func DeRegister() {
	getDiscoverClient().DeRegister()
}
