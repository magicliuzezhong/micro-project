//
// Package common
// @Description：从注册中心获取到的实例
// @Author：liuzezhong 2021/6/28 10:02 上午
// @Company cloud-ark.com
//
package common

import "strconv"

//
// ServiceInstance
// @Description: 服务实例
//
type ServiceInstance struct {
	// 主机
	Host string
	// 端口
	Port int
	// 权重
	Weight int
	// 当前权重
	CurWeight int
	// grpc端口
	GrpcPort int
	// 调用次数
	CallTimes int
}

//
// GetUrl
// @Description: 获取url
// @receiver instance 实例
// @return string url
//
func (instance ServiceInstance) GetUrl() string {
	var url = ""
	if instance.Host != "" && instance.Port != 0 {
		url += "http://" + instance.Host + ":" + strconv.Itoa(instance.Port)
	}
	return url
}
