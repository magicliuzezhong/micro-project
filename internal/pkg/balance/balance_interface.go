//
// Package balance
// @Description：负载均衡接口
// @Author：liuzezhong 2021/6/30 5:20 下午
// @Company cloud-ark.com
//
package balance

import "micro-project/internal/pkg/common"

//
// ILoadBalance
// @Description: 负载均衡接口
//
type ILoadBalance interface {
	//
	// DoBalance
	// @Description: 负载均衡
	// @param instance 待均衡实例
	// @return common.ServiceInstance 实例
	// @return error 错误
	//
	DoBalance([]*common.ServiceInstance, ...string) (*common.ServiceInstance, error)
}
