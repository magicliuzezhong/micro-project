//
// Package balance
// @Description：随机负载均衡
// @Author：liuzezhong 2021/6/30 5:23 下午
// @Company cloud-ark.com
//
package balance

import (
	"errors"
	"math/rand"
	"micro-project/internal/pkg/common"
	"sync"
)

//
// RandomBalance
// @Description: 随机负载均衡实现
//
type randomBalance struct {
}

//
// @Description: 用于单例模式
//
var randomBalanceOnce sync.Once

//
// @Description: 实例
//
var randomBalanceInstance *randomBalance

//
// NewRandomBalance
// @Description: 获取随机负载均衡实例
// @return ILoadBalance 负载均衡接口
//
func NewRandomBalance() ILoadBalance {
	randomBalanceOnce.Do(func() {
		randomBalanceInstance = &randomBalance{}
	})
	return randomBalanceInstance
}

//
// DoBalance
// @Description: 随机负载均衡
// @receiver r RandomLoadBalance
// @param instances 实例
// @return *common.ServiceInstance 返回实例
// @return error 错误
//
func (r randomBalance) DoBalance(instances []*common.ServiceInstance, _ ...string) (*common.ServiceInstance, error) {
	lens := len(instances)
	if lens == 0 {
		return nil, errors.New("实例找不到")
	}
	index := rand.Intn(lens)
	instance := instances[index]
	instance.CallTimes++
	return instance, nil
}
