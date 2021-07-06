//
// Package balance
// @Description：轮询负载均衡实现
// @Author：liuzezhong 2021/6/30 5:31 下午
// @Company cloud-ark.com
//
package balance

import (
	"fmt"
	"micro-project/internal/pkg/common"
	"sync"
)

//
// roundRobinBalance
// @Description: 轮询
//
type roundRobinBalance struct {
	currentIndex int
}

//
// @Description: 单例
//
var roundRobinBalanceOnce sync.Once

//
// @Description: 实例
//
var roundRobinBalanceInstance *roundRobinBalance

//
// NewRoundRobinBalance
// @Description: 获取轮询实例
// @return ILoadBalance 负载均衡接口
//
func NewRoundRobinBalance() ILoadBalance {
	return &roundRobinBalance{
		currentIndex: 0,
	}
}

//
// DoBalance
// @Description: 轮询
// @receiver r roundRobinBalance
// @param instances 实例
// @return *common.ServiceInstance 返回实例
// @return error 错误
//
func (r *roundRobinBalance) DoBalance(instances []*common.ServiceInstance,
	_ ...string) (*common.ServiceInstance, error) {
	lens := len(instances)
	if lens == 0 {
		return nil, fmt.Errorf("实例找不到")
	}
	if r.currentIndex >= lens {
		r.currentIndex = 0
	}
	var instance = instances[r.currentIndex]
	r.currentIndex = (r.currentIndex + 1) % lens
	instance.CallTimes++
	return instance, nil
}
