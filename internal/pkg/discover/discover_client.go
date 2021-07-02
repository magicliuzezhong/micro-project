//
// Package discover
// @Description：服务注册通用接口
// @Author：liuzezhong 2021/6/28 9:54 上午
// @Company cloud-ark.com
//
package discover

import (
	"micro-project/internal/pkg/common"
)

//
// IDiscoveryClient
// @Description: 服务注册通用接口
//
type IDiscoveryClient interface {
	//
	// Register
	// @Description: 服务注册接口
	// @param svcHost 注册host
	// @param svcPort 注册port
	// @param svcName 注册服务名称
	// @param weight 权重
	// @param meta 元数据
	// @param tags 标签
	// @return bool 是否注册成功
	//
	Register(svcHost string, svcPort int, svcName string, weight int, meta map[string]string, tags []string) bool
	//
	// DeRegister
	// @Description: 服务注销接口
	// @return bool 注销成功
	//
	DeRegister() bool
	//
	// DiscoverServices
	// @Description: 发现服务实例接口
	// @param serviceName 服务名
	// @return []*common.ServiceInstance 服务实例
	//
	DiscoverServices(serviceName string) []*common.ServiceInstance
}
